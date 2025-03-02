package cmd

import (
	"fmt"
	"math/rand/v2"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/daos"
	"github.com/spf13/cobra"

	"github.com/fragoulis/setip_v2/internal/db/company"
	"github.com/fragoulis/setip_v2/internal/db/employment"
	"github.com/fragoulis/setip_v2/internal/db/member"
	"github.com/fragoulis/setip_v2/internal/db/payment"
	"github.com/fragoulis/setip_v2/internal/db/subscription"
)

const (
	Percent         = 100
	DefaultBranchPc = 30
)

func NewSeedCommand(app core.App) *cobra.Command {
	//nolint:errcheck
	gofakeit.Seed(0)

	command := &cobra.Command{
		Use:   "seed",
		Short: "Generates random data",
	}

	command.AddCommand(seedCompanyCommand(app))
	command.AddCommand(seedMemberCommand(app))
	command.AddCommand(seedPaymentCommand(app))

	return command
}

func seedCompanyCommand(app core.App) *cobra.Command {
	branchPc := 0
	save := false
	count := 0

	command := &cobra.Command{
		Use:          "company",
		Example:      "seed company",
		Short:        "Generates a random company",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			for i := 0; i < count; i++ {
				err := app.Dao().RunInTransaction(func(tx *daos.Dao) error {
					// Set the branch for a subset of the generated companies.
					withBranch := rand.IntN(Percent) < branchPc

					company, err := company.NewRandom(tx, withBranch)
					if err != nil {
						return fmt.Errorf("failed to fetch random company: %w", err)
					}

					fmt.Printf("%+v\n", company)

					if !save {
						return nil
					}

					return tx.Save(company)
				})
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	command.Flags().IntVar(&branchPc, "branch-pc", DefaultBranchPc, "percentage with branch")
	command.Flags().BoolVar(&save, "save", false, "save generated records to database")
	command.Flags().IntVar(&count, "count", 1, "number of records to generate")

	return command
}

func seedMemberCommand(app core.App) *cobra.Command {
	var (
		save         bool
		count        int
		unpaidFeePc  int
		inactivePc   int
		unemployedPc int
	)

	command := &cobra.Command{
		Use:          "member",
		Example:      "seed member",
		Short:        "Generates a random member",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			for i := 0; i < count; i++ {
				err := app.Dao().RunInTransaction(func(tx *daos.Dao) error {
					// 1. Create member
					// 2. Create subscription
					// 3. Create employment

					member, err := member.NewRandom(tx)
					if err != nil {
						return fmt.Errorf("random member: %w", err)
					}

					// Generate the ID here so that we can use it later even
					// without creating the record.
					member.RefreshId()

					if save {
						err = tx.Save(member)
						if err != nil {
							return err
						}
					}

					// Set fee paid and active randomly based on the percetages
					// passed via the flags.
					paid := rand.IntN(100) > unpaidFeePc
					active := rand.IntN(100) > inactivePc

					subscription, err := subscription.NewRandom(
						member.GetId(),
						paid,
						active,
					)
					if err != nil {
						return fmt.Errorf("random subscription: %w", err)
					}

					if save {
						err = tx.Save(subscription)
						if err != nil {
							return err
						}
					}

					employed := rand.IntN(100) > unemployedPc

					employment, err := employment.NewRandom(
						tx,
						member.GetId(),
						employed,
					)
					if err != nil {
						return fmt.Errorf("random employment: %w", err)
					}

					if save {
						err = tx.Save(employment)
						if err != nil {
							return err
						}
					}

					fmt.Printf("%+v\n", member)
					fmt.Printf("%+v\n", subscription)
					fmt.Printf("%+v\n", employment)

					return nil
				})
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	command.Flags().IntVar(&unpaidFeePc, "fee-paid-pc", 10, "percentage with unpaid fee")
	command.Flags().IntVar(&inactivePc, "active-pc", 10, "percentage inactive")
	command.Flags().IntVar(&unemployedPc, "unemployed-pc", 5, "percentage unemployed")
	command.Flags().BoolVar(&save, "save", false, "save generated records to database")
	command.Flags().IntVar(&count, "count", 1, "number of records to generate")

	return command
}

func seedPaymentCommand(app core.App) *cobra.Command {
	var (
		save  bool
		count int
	)

	command := &cobra.Command{
		Use:          "payment",
		Example:      "seed payment",
		Short:        "Generates a random payment",
		SilenceUsage: true,
		RunE: func(_ *cobra.Command, _ []string) error {
			for i := 0; i < count; i++ {
				err := app.Dao().RunInTransaction(func(tx *daos.Dao) error {
					payment, err := payment.NewRandom(tx)
					if err != nil {
						return fmt.Errorf("random payment: %w", err)
					}

					fmt.Printf("%+v\n", payment)

					if !save {
						return nil
					}

					return tx.Save(payment)
				})
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	command.Flags().BoolVar(&save, "save", false, "save generated records to database")
	command.Flags().IntVar(&count, "count", 1, "number of records to generate")

	return command
}
