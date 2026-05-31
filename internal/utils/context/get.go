package context

import (
	"context"
	"errors"
	"go-learning/internal/utils/cache"
	"log"
)

type InfoUserUUID struct {
	UserId      int32
	UserAccount string
}

func GetSubjectUUID(ctx context.Context) (string, error) {
	sUUID, ok := ctx.Value("subjectUUID").(string)
	if !ok {
		return "", errors.New("failed to get subject UUID")
	}

	return sUUID, nil
}

func GetUserIdFromUUID(ctx context.Context) (int32, error) {
	sUUID, err := GetSubjectUUID(ctx)
	log.Println("sUUID: ", sUUID)
	if err != nil {
		return 0, err
	}

	// get infoUser Redis from uuid
	var infoUser InfoUserUUID
	if err = cache.GetCache(ctx, sUUID, &infoUser); err != nil {
		return 0, err
	}

	return infoUser.UserId, nil
}
