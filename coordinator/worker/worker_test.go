package worker

import (
	"testing"
	"time"

	"github.com/bgentry/que-go"
	gm "github.com/onsi/gomega"

	"github.com/capeprivacy/cape/coordinator"
	"github.com/capeprivacy/cape/coordinator/database"
	"github.com/capeprivacy/cape/framework"
	errors "github.com/capeprivacy/cape/partyerrors"
	"github.com/capeprivacy/cape/primitives"
)

func TestWorkerDeleteRecoveries(t *testing.T) {
	gm.RegisterTestingT(t)

	url, err := primitives.NewURL("https://www.test.com")
	gm.Expect(err).To(gm.BeNil())

	t.Run("no recoveries found", func(t *testing.T) {
		mock, err := coordinator.NewMockClientTransport(url, []*coordinator.MockResponse{})
		gm.Expect(err).To(gm.BeNil())

		w := &Worker{
			coordClient: coordinator.NewClient(mock),
			logger:      framework.TestLogger(),
		}

		err = w.DeleteRecoveries(&que.Job{})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(mock.Requests)).To(gm.Equal(1))
		gm.Expect(mock.Requests[0].Query).To(gm.ContainSubstring("ListRecoveries"))
	})

	t.Run("no expired recoveries", func(t *testing.T) {
		recoveryA, err := primitives.GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())

		recoveryB, err := primitives.GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())

		mock, err := coordinator.NewMockClientTransport(url, []*coordinator.MockResponse{
			{
				Value: coordinator.ListRecoveriesResponse{
					Recoveries: []*primitives.Recovery{
						recoveryA,
						recoveryB,
					},
				},
			},
		})
		gm.Expect(err).To(gm.BeNil())

		w := &Worker{
			coordClient: coordinator.NewClient(mock),
			logger:      framework.TestLogger(),
		}

		err = w.DeleteRecoveries(&que.Job{})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(mock.Requests)).To(gm.Equal(1))
		gm.Expect(mock.Requests[0].Query).To(gm.ContainSubstring("ListRecoveries"))
	})

	t.Run("deletes only expired recoveries", func(t *testing.T) {
		recoveryA, err := primitives.GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())

		recoveryB, err := primitives.GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())

		recoveryC, err := primitives.GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())

		recoveryC.ExpiresAt = time.Now().UTC().Add(-1 * time.Minute)

		mock, err := coordinator.NewMockClientTransport(url, []*coordinator.MockResponse{
			{
				Value: coordinator.ListRecoveriesResponse{
					Recoveries: []*primitives.Recovery{
						recoveryA,
						recoveryB,
						recoveryC,
					},
				},
			},
			{
				Value: nil,
			},
		})
		gm.Expect(err).To(gm.BeNil())

		w := &Worker{
			coordClient: coordinator.NewClient(mock),
			logger:      framework.TestLogger(),
		}

		err = w.DeleteRecoveries(&que.Job{})
		gm.Expect(err).To(gm.BeNil())

		gm.Expect(len(mock.Requests)).To(gm.Equal(2))
		gm.Expect(mock.Requests[0].Query).To(gm.ContainSubstring("ListRecoveries"))
		gm.Expect(mock.Requests[1].Query).To(gm.ContainSubstring("DeleteRecoveries"))

		ids := mock.Requests[1].Variables["ids"].([]database.ID)
		gm.Expect(len(ids)).To(gm.Equal(1))
		gm.Expect(ids[0]).To(gm.Equal(recoveryC.ID))
	})

	t.Run("list recoveries errors", func(t *testing.T) {
		mock, err := coordinator.NewMockClientTransport(url, []*coordinator.MockResponse{
			{
				Error: errors.ErrNotImplemented,
			},
		})
		gm.Expect(err).To(gm.BeNil())

		w := &Worker{
			coordClient: coordinator.NewClient(mock),
			logger:      framework.TestLogger(),
		}

		err = w.DeleteRecoveries(&que.Job{})
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.ContainSubstring("not_implemented"))
	})

	t.Run("delete recoveries errors", func(t *testing.T) {
		recoveryA, err := primitives.GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())

		recoveryB, err := primitives.GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())

		recoveryC, err := primitives.GenerateRecovery()
		gm.Expect(err).To(gm.BeNil())

		recoveryC.ExpiresAt = time.Now().UTC().Add(-1 * time.Minute)

		mock, err := coordinator.NewMockClientTransport(url, []*coordinator.MockResponse{
			{
				Value: coordinator.ListRecoveriesResponse{
					Recoveries: []*primitives.Recovery{
						recoveryA,
						recoveryB,
						recoveryC,
					},
				},
			},
			{
				Error: errors.ErrNotImplemented,
			},
		})
		gm.Expect(err).To(gm.BeNil())

		w := &Worker{
			coordClient: coordinator.NewClient(mock),
			logger:      framework.TestLogger(),
		}

		err = w.DeleteRecoveries(&que.Job{})
		gm.Expect(err).ToNot(gm.BeNil())
		gm.Expect(err.Error()).To(gm.ContainSubstring("not_implemented"))

		gm.Expect(len(mock.Requests)).To(gm.Equal(2))
		gm.Expect(mock.Requests[0].Query).To(gm.ContainSubstring("ListRecoveries"))
		gm.Expect(mock.Requests[1].Query).To(gm.ContainSubstring("DeleteRecoveries"))
	})
}
