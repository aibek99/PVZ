package create

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"Homework-1/internal/app/command"
	"Homework-1/internal/model/cli"
	"Homework-1/internal/storage/pvz"
	"Homework-1/internal/storage/pvz/mock"
	"Homework-1/pkg/constants"
	"Homework-1/pkg/errlst"
)

func TestCommand_Do(t *testing.T) {
	t.Parallel()
	msg := constants.SuccessfullyCreatedPVZ
	ctrl := minimock.NewController(t)
	tests := []struct {
		description string
		args        *command.Arguments
		storage     pvz.Storage
		want        *string
		wantErr     error
	}{
		{
			description: "",
			args: &command.Arguments{
				Name:    "Sample PVZ",
				Address: "kazan",
				Contact: "+4566545454",
			},
			storage: mock.NewStorageMock(ctrl).
				CreateMock.When(
				&cli.PVZ{
					Name:    "Sample PVZ",
					Address: "kazan",
					Contact: "+4566545454",
				}).
				Then(nil),
			wantErr: nil,
			want:    &msg,
		},
		{
			description: "",
			args: &command.Arguments{
				Name:    "Sample PVZ",
				Address: "kazan",
				Contact: "+4566545454",
			},
			storage: mock.NewStorageMock(ctrl).
				CreateMock.When(
				&cli.PVZ{
					Name:    "Sample PVZ",
					Address: "kazan",
					Contact: "+4566545454",
				}).
				Then(errlst.ErrPVZAlreadyExists),
			wantErr: errlst.ErrPVZAlreadyExists,
			want:    nil,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.description, func(t *testing.T) {
			t.Parallel()
			pvzCommand := New(tt.storage)

			message, err := pvzCommand.Do(tt.args)

			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, message)
		})
	}
}
