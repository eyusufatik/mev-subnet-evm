// (c) 2019-2020, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package evm

import (
	"fmt"
	"net/http"

	"github.com/ava-labs/avalanchego/api"
	"github.com/ava-labs/avalanchego/utils/profiler"
	"github.com/ethereum/go-ethereum/log"
)

// Admin is the API service for admin API calls
type Admin struct {
	vm       *VM
	profiler profiler.Profiler
}

func NewAdminService(vm *VM, performanceDir string) *Admin {
	return &Admin{
		vm:       vm,
		profiler: profiler.New(performanceDir),
	}
}

// StartCPUProfiler starts a cpu profile writing to the specified file
func (p *Admin) StartCPUProfiler(r *http.Request, args *struct{}, reply *api.EmptyReply) error {
	log.Info("Admin: StartCPUProfiler called")

	err := p.profiler.StartCPUProfiler()
	return err
}

// StopCPUProfiler stops the cpu profile
func (p *Admin) StopCPUProfiler(r *http.Request, args *struct{}, reply *api.EmptyReply) error {
	log.Info("Admin: StopCPUProfiler called")

	err := p.profiler.StopCPUProfiler()
	return err
}

// MemoryProfile runs a memory profile writing to the specified file
func (p *Admin) MemoryProfile(r *http.Request, args *struct{}, reply *api.EmptyReply) error {
	log.Info("Admin: MemoryProfile called")

	err := p.profiler.MemoryProfile()
	return err
}

// LockProfile runs a mutex profile writing to the specified file
func (p *Admin) LockProfile(r *http.Request, args *struct{}, reply *api.EmptyReply) error {
	log.Info("Admin: LockProfile called")

	err := p.profiler.LockProfile()
	return err
}

type SetLogLevelArgs struct {
	Level string `json:"level"`
}

func (p *Admin) SetLogLevel(r *http.Request, args *SetLogLevelArgs, reply *api.EmptyReply) error {
	log.Info("EVM: SetLogLevel called", "logLevel", args.Level)
	if err := p.vm.logger.SetLogLevel(args.Level); err != nil {
		return fmt.Errorf("failed to parse log level: %w ", err)
	}
	return nil
}

type ConfigReply struct {
	Config *Config `json:"config"`
}

func (p *Admin) GetVMConfig(r *http.Request, args *struct{}, reply *ConfigReply) error {
	reply.Config = &p.vm.config
	return nil
}
