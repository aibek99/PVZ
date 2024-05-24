package get

import (
	"testing"

	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"

	"Homework-1/internal/app/command"
	"Homework-1/internal/model/cli"
	"Homework-1/internal/storage/pvz"
	"Homework-1/internal/storage/pvz/mock"
	"Homework-1/pkg/errlst"
)

func TestCommand_Do(t *testing.T) {
	t.Parallel()
	msg := "\nPVZ Name: Sample PVZ\nPVZ address: kazan\nPVZ Contact: +4566545454\n"
	ctrl := minimock.NewController(t)
	tests := []struct {
		description string
		args        *command.Arguments
		storage     pvz.Storage
		want        *string
		wantErr     error
	}{
		{
			description: "PVZ was found",
			args: &command.Arguments{
				Name: "Sample PVZ",
			},
			storage: mock.NewStorageMock(ctrl).FindMock.When("Sample PVZ").Then(
				cli.PVZ{
					Name:    "Sample PVZ",
					Address: "kazan",
					Contact: "+4566545454",
				}, nil,
			),
			want:    &msg,
			wantErr: nil,
		},
		{
			description: "PVZ was not found",
			args: &command.Arguments{
				Name: "Sample PVZ",
			},
			storage: mock.NewStorageMock(ctrl).FindMock.When("Sample PVZ").Then(cli.PVZ{}, errlst.ErrPVZNotFound),
			wantErr: errlst.ErrPVZNotFound,
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
