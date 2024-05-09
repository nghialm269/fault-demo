package users

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Southclaws/fault/fmsg"
	"github.com/Southclaws/fault/ftag"

	"github.com/nghialm269/fault-demo/pkg/fserrors"
	"github.com/nghialm269/fault-demo/pkg/fserrors/dberrors"
	"github.com/nghialm269/fault-demo/pkg/fserrors/wrappers/errctx"
)

func HandlerGetUserByID(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	rawID := r.PathValue("id")
	id, err := strconv.ParseInt(rawID, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return err
	}

	user, err := ServiceGetUser(ctx, ServiceGetUserParams{ID: id})
	if err != nil {

		var status int
		switch ftag.Get(err) {
		case ftag.NotFound:
			status = http.StatusNotFound
		case ftag.InvalidArgument:
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)

		json.NewEncoder(w).Encode(struct {
			Message string `json:"message"`
		}{
			Message: fmsg.GetIssue(err),
		})

		return fserrors.Wrap(err,
			errctx.With(ctx, "id", id),
			fmsg.With("ServiceGetUser"),
		)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(user)
}

func ServiceGetUser(ctx context.Context, params ServiceGetUserParams) (ServiceGetUserResult, error) {
	result, err := RepositoryGetUserByID(ctx, RepositoryGetUserByIDParams{
		ID: params.ID,
	})
	if err != nil {
		return ServiceGetUserResult{}, fserrors.Wrap(err,
			errctx.With(ctx, "serviceParams", params),
			fmsg.WithDesc("RepositoryGetUserByID", "Failed to get user, please try again!"),
		)
	}

	return ServiceGetUserResult{
		ID:       result.ID,
		Username: result.Username,
	}, nil
}

func RepositoryGetUserByID(ctx context.Context, params RepositoryGetUserByIDParams) (RepositoryGetUserByIDResult, error) {
	ctx = errctx.WithMeta(ctx, "repositoryParams", params)

	switch {
	case params.ID == 404:
		return RepositoryGetUserByIDResult{}, fserrors.Wrap(dberrors.ErrEntryNotFound,
			errctx.With(ctx),
			ftag.With(ftag.NotFound),
		)
	case params.ID >= 400 && params.ID < 500:
		message := fmt.Sprintf("error %d", params.ID)
		return RepositoryGetUserByIDResult{}, fserrors.New(message,
			errctx.With(ctx),
			ftag.With(ftag.InvalidArgument),
		)

	case params.ID >= 500 && params.ID < 600:
		message := fmt.Sprintf("error %d", params.ID)
		return RepositoryGetUserByIDResult{}, fserrors.New(message,
			errctx.With(ctx),
			ftag.With(ftag.Internal),
		)

	case params.ID == 999:
		panic("panic 999")
	}

	return RepositoryGetUserByIDResult{
		ID:       params.ID,
		Username: fmt.Sprintf("user%d@example.com", params.ID),
	}, nil
}
