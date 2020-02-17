// Copyright 2019-present Open Networking Foundation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package _map

import (
	"context"
	"github.com/onosproject/onos-test/pkg/input"
	"github.com/onosproject/onos-test/pkg/simulation"
	"time"

	"github.com/atomix/go-client/pkg/client/map"
	"github.com/onosproject/onos-test/pkg/onit/env"
	"github.com/onosproject/onos-test/pkg/onit/setup"
)

const (
	keyLengthArg   = "key-length"
	keyCountArg    = "key-count"
	valueLengthArg = "value-length"
)

const (
	defaultKeyLength   = 8
	defaultKeyCount    = 10
	defaultValueLength = 128
)

// MapSimulation :: simulation
type MapSimulation struct {
	simulation.Suite
	_map    _map.Map
	watchCh chan *_map.Event
	keys    input.Source
	values  input.Source
}

// SetupSimulation :: simulation
func (s *MapSimulation) SetupSimulation(sim *simulation.Simulator) error {
	setup.Database("raft").Raft()
	return setup.Setup()
}

// ScheduleSimulator :: simulation
func (s *MapSimulation) ScheduleSimulator(sim *simulation.Simulator) {
	sim.Schedule("put", s.SimulateMapPut, 1*time.Second, 1)
	sim.Schedule("get", s.SimulateMapGet, 1*time.Second, 1)
	sim.Schedule("remove", s.SimulateMapRemove, 1*time.Second, 1)
	sim.Schedule("event", s.SimulateMapEvent, 5*time.Second, 1)
}

// SetupSimulator :: simulation
func (s *MapSimulation) SetupSimulator(sim *simulation.Simulator) error {
	database, err := env.Storage().Database("raft").Connect()
	if err != nil {
		panic(err)
	}
	m, err := database.GetMap(context.Background(), sim.Name)
	if err != nil {
		panic(err)
	}
	s._map = m

	watchCh := make(chan *_map.Event)
	if err := s._map.Watch(context.Background(), watchCh); err != nil {
		return err
	}
	s.watchCh = watchCh
	s.setupInputs(sim)
	return nil
}

// setupInputs sets up the simulator inputs
func (s *MapSimulation) setupInputs(sim *simulation.Simulator) {
	s.keys = input.RandomChoice(
		input.SetOf(
			input.RandomString(
				sim.Arg(keyLengthArg).Int(defaultKeyLength)),
			sim.Arg(keyCountArg).Int(defaultKeyCount)))

	s.values = input.RandomBytes(sim.Arg(valueLengthArg).Int(defaultValueLength))
}

// nextKey returns the next simulator key
func (s *MapSimulation) nextKey(sim *simulation.Simulator) string {
	return s.keys.Next().String()
}

// nextValue returns the next simulator value
func (s *MapSimulation) nextValue(sim *simulation.Simulator) string {
	return string(s.values.Next().Bytes())
}

// TearDownSimulator :: simulation
func (s *MapSimulation) TearDownSimulator(c *simulation.Simulator) error {
	_ = s._map.Close(context.Background())
	return nil
}

// SimulateMapPut :: simulation
func (s *MapSimulation) SimulateMapPut(sim *simulation.Simulator) error {
	kv, err := s._map.Put(context.Background(), s.nextKey(sim), []byte(s.nextValue(sim)))
	if err != nil {
		return err
	}
	sim.TraceFields("op", "put", "process", sim.Process, "key", kv.Key, "value", string(kv.Value), "version", kv.Version)
	return nil
}

// SimulateMapGet :: simulation
func (s *MapSimulation) SimulateMapGet(sim *simulation.Simulator) error {
	key := s.nextKey(sim)
	kv, err := s._map.Get(context.Background(), key)
	if err != nil {
		return err
	}
	if kv != nil {
		sim.TraceFields("op", "get", "process", sim.Process, "key", kv.Key, "value", string(kv.Value), "version", kv.Version)
	}
	return nil
}

// SimulateMapRemove :: simulation
func (s *MapSimulation) SimulateMapRemove(sim *simulation.Simulator) error {
	key := s.nextKey(sim)
	kv, err := s._map.Get(context.Background(), key)
	if err != nil {
		return err
	}
	if kv != nil {
		sim.TraceFields("op", "remove", "process", sim.Process, "key", kv.Key, "value", string(kv.Value), "version", kv.Version)
	}
	return nil
}

// SimulateMapEvent :: simulation
func (s *MapSimulation) SimulateMapEvent(sim *simulation.Simulator) error {
	select {
	case event := <-s.watchCh:
		sim.TraceFields("op", "event", "process", sim.Process, "key", event.Entry.Key, "value", string(event.Entry.Value), "version", event.Entry.Version)
		return nil
	case <-time.After(5 * time.Second):
		return nil
	}
}
