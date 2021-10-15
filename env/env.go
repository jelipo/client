package env

import "strings"

type DiffEnvs struct {
	NewEnv     []InternalEnv
	DeletedEnv []InternalEnv
	ChangedEnv []InternalEnv
}

func Compare(beforeEnvs []InternalEnv, afterEnvs []InternalEnv) DiffEnvs {
	beforeEnvsMap := sliceToMap(beforeEnvs)
	afterEnvsMap := sliceToMap(afterEnvs)
	var deletedEnvs []InternalEnv
	var changedEnvs []InternalEnv
	var newEnvs []InternalEnv

	for beforeName, beforeEnv := range beforeEnvsMap {
		afterEnv := afterEnvsMap[beforeName]
		if afterEnv == nil {
			deletedEnvs = append(deletedEnvs, *beforeEnv)
		} else if afterEnv.Value != beforeEnv.Value {
			changedEnvs = append(changedEnvs, *afterEnv)
		}
	}

	for afterName, afterEnv := range afterEnvsMap {
		beforeEnv := beforeEnvsMap[afterName]
		if beforeEnv == nil {
			newEnvs = append(newEnvs, *afterEnv)
		}
	}
	return DiffEnvs{
		NewEnv:     newEnvs,
		DeletedEnv: deletedEnvs,
		ChangedEnv: changedEnvs,
	}
}

type InternalEnv struct {
	Name      string
	Value     string
	SecureEnv bool
}

func ReadEnvs(envs []string) []InternalEnv {
	var internalEnvs []InternalEnv
	for _, env := range envs {
		split := strings.Split(env, "=")
		if len(split) != 2 {
			continue
		}
		internalEnv := InternalEnv{
			Name:      split[0],
			Value:     split[1],
			SecureEnv: false,
		}
		internalEnvs = append(internalEnvs, internalEnv)
	}
	return internalEnvs
}

func sliceToMap(envs []InternalEnv) map[string]*InternalEnv {
	mapObj := map[string]*InternalEnv{}
	for _, env := range envs {
		mapObj[env.Name] = &env
	}
	return mapObj
}
