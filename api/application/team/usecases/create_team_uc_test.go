package usecases

import (
	"reflect"
	"sportlink/api/domain/player"
	"sportlink/api/domain/team"
	"testing"
)

func TestCreateTeamUC_Invoke(t *testing.T) {
	type fields struct {
		playerRepository player.Repository
		teamRepository   team.Repository
	}
	type args struct {
		input team.Entity
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *team.Entity
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := &CreateTeamUC{
				playerRepository: tt.fields.playerRepository,
				teamRepository:   tt.fields.teamRepository,
			}
			got, err := uc.Invoke(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Invoke() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Invoke() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewCreateTeamUC(t *testing.T) {
	type args struct {
		playerRepository player.Repository
		teamRepository   team.Repository
	}
	tests := []struct {
		name string
		args args
		want *CreateTeamUC
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCreateTeamUC(tt.args.playerRepository, tt.args.teamRepository); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCreateTeamUC() = %v, want %v", got, tt.want)
			}
		})
	}
}
