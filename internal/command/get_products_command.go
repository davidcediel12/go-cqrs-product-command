package command

type GetProductsCommand struct {
	Page             int
	Size             int
	Images           bool
	OnlyPrimaryImage bool
}
