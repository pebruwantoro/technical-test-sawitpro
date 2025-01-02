// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import "context"

type RepositoryInterface interface {
	CreateEstate(ctx context.Context, input Estate) (result Estate, err error)
	CreateEstateTree(ctx context.Context, input EstateTree) (result EstateTree, err error)
	GetStatsByEstateId(ctx context.Context, id string) (result StatsEstate, err error)
	GetEstateById(ctx context.Context, id string) (result Estate, err error)
	GetTreesByEstateId(ctx context.Context, id string) (result []EstateTree, err error)
}
