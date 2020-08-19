package commit

// func newCommitBuildCmd(container *lib.Container) *cobra.Command {
// 	executor := container.Executor
// 	return &cobra.Command{
// 		Use:   "build",
// 		Short: "Commits the build catalog for the given service",
// 		Args: func(cmd *cobra.Command, args []string) error {
// 			if len(args) < 1 {
// 				return errors.New(styles.StyleError("Argument <service> required"))
// 			}

// 			if len(args) > 1 {
// 				return errors.New(styles.StyleError("Unexpected arguments, only <service> expected"))
// 			}

// 			fs := lib.CreateOsFs()
// 			context, err := lib.GetContext(fs, cmd)
// 			if err != nil {
// 				panic(err)
// 			}

// 			var serviceNames []string

// 			// for serviceName := range context.Config.Services {
// 			// 	if serviceName == args[0] {
// 			// 		return nil
// 			// 	}
// 			// 	serviceNames = append(serviceNames, serviceName)
// 			// }

// 			return errors.New(styles.StyleError(`Unknown service "` + args[0] + `". Valid options: ` + styles.StyleEnumOptions(serviceNames) + "."))
// 		},
// 		Run: func(cmd *cobra.Command, args []string) {
// 			service := args[0]
// 			fs := lib.CreateOsFs()
// 			context, err := lib.GetContext(fs, cmd)
// 			if err != nil {
// 				panic(err)
// 			}

// 			// serviceManifest, err := context.LoadServiceManifest(service)
// 			// if err != nil {
// 			// 	panic(err)
// 			// }

// 			//executor.ExecCommand(fmt.Sprintf("git add .g3ops/builds/%s.yml", service), lib.ExecOptions{DryRun: context.DryRun})
// 			//executor.ExecCommand(fmt.Sprintf(`git commit -m "Generated build for %s %s"`, service, serviceManifest.Version), lib.ExecOptions{DryRun: context.DryRun})
// 		},
// 	}
// }
