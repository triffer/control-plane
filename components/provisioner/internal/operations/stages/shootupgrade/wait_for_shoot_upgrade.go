package shootupgrade

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kyma-project/control-plane/components/provisioner/internal/provisioning/persistence/dbsession"

	gardener_typeshelper "github.com/gardener/gardener/pkg/apis/core/v1beta1/helper"

	gardener_types "github.com/gardener/gardener/pkg/apis/core/v1beta1"
	"github.com/kyma-project/control-plane/components/provisioner/internal/model"
	"github.com/kyma-project/control-plane/components/provisioner/internal/operations"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GardenerClient interface {
	Get(ctx context.Context, name string, options v1.GetOptions) (*gardener_types.Shoot, error)
}

//go:generate mockery --name=KubeconfigProvider
type KubeconfigProvider interface {
	FetchRaw(context.Context, gardener_types.Shoot) ([]byte, error)
}

type WaitForShootUpgradeStep struct {
	gardenerClient GardenerClient
	nextStep       model.OperationStage
	timeLimit      time.Duration

	dbSession          dbsession.ReadWriteSession
	kubeconfigProvider KubeconfigProvider
}

func NewWaitForShootUpgradeStep(
	gardenerClient GardenerClient,
	dbSession dbsession.ReadWriteSession,
	kubeconfigProvider KubeconfigProvider,
	nextStep model.OperationStage,
	timeLimit time.Duration) *WaitForShootUpgradeStep {

	return &WaitForShootUpgradeStep{
		gardenerClient: gardenerClient,
		nextStep:       nextStep,
		timeLimit:      timeLimit,

		dbSession:          dbSession,
		kubeconfigProvider: kubeconfigProvider,
	}
}

func (s WaitForShootUpgradeStep) Name() model.OperationStage {
	return model.WaitingForShootUpgrade
}

func (s *WaitForShootUpgradeStep) TimeLimit() time.Duration {
	return s.timeLimit
}

func (s *WaitForShootUpgradeStep) Run(cluster model.Cluster, _ model.Operation, logger logrus.FieldLogger) (operations.StageResult, error) {

	gardenerConfig := cluster.ClusterConfig

	shoot, err := s.gardenerClient.Get(context.Background(), gardenerConfig.Name, v1.GetOptions{})
	if err != nil {
		return operations.StageResult{}, err
	}

	lastOperation := shoot.Status.LastOperation

	if lastOperation != nil {
		if lastOperation.State == gardener_types.LastOperationStateSucceeded {

			kubeconfig, err := s.kubeconfigProvider.FetchRaw(context.TODO(), *shoot)
			if err != nil {
				return operations.StageResult{}, err
			}
			if dberr := s.dbSession.UpdateKubeconfig(cluster.ID, string(kubeconfig)); dberr != nil {
				return operations.StageResult{}, dberr
			}

			return operations.StageResult{Stage: s.nextStep, Delay: 0}, nil
		}

		if lastOperation.State == gardener_types.LastOperationStateFailed {
			if gardener_typeshelper.HasErrorCode(shoot.Status.LastErrors, gardener_types.ErrorInfraRateLimitsExceeded) {
				return operations.StageResult{}, errors.New("error during shoot cluster upgrade: rate limits exceeded")
			}
			logger.Warningf("Gardener Shoot cluster upgrade operation failed! Last state: %s, Description: %s", lastOperation.State, lastOperation.Description)

			err := fmt.Errorf("gardener Shoot cluster upgrade failed. Last Shoot state: %s, Shoot description: %s", lastOperation.State, lastOperation.Description)
			return operations.StageResult{}, operations.NewNonRecoverableError(err)
		}
	}

	return operations.StageResult{Stage: s.Name(), Delay: 20 * time.Second}, nil
}
