// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/openshift/api/machine/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// MachineStatusApplyConfiguration represents an declarative configuration of the MachineStatus type for use
// with apply.
type MachineStatusApplyConfiguration struct {
	NodeRef                *v1.ObjectReference              `json:"nodeRef,omitempty"`
	LastUpdated            *metav1.Time                     `json:"lastUpdated,omitempty"`
	ErrorReason            *v1beta1.MachineStatusError      `json:"errorReason,omitempty"`
	ErrorMessage           *string                          `json:"errorMessage,omitempty"`
	ProviderStatus         *runtime.RawExtension            `json:"providerStatus,omitempty"`
	Addresses              []v1.NodeAddress                 `json:"addresses,omitempty"`
	LastOperation          *LastOperationApplyConfiguration `json:"lastOperation,omitempty"`
	Phase                  *string                          `json:"phase,omitempty"`
	Conditions             []ConditionApplyConfiguration    `json:"conditions,omitempty"`
	AuthoritativeAPI       *v1beta1.MachineAuthority        `json:"authoritativeAPI,omitempty"`
	SynchronizedGeneration *int64                           `json:"synchronizedGeneration,omitempty"`
}

// MachineStatusApplyConfiguration constructs an declarative configuration of the MachineStatus type for use with
// apply.
func MachineStatus() *MachineStatusApplyConfiguration {
	return &MachineStatusApplyConfiguration{}
}

// WithNodeRef sets the NodeRef field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the NodeRef field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithNodeRef(value v1.ObjectReference) *MachineStatusApplyConfiguration {
	b.NodeRef = &value
	return b
}

// WithLastUpdated sets the LastUpdated field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastUpdated field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithLastUpdated(value metav1.Time) *MachineStatusApplyConfiguration {
	b.LastUpdated = &value
	return b
}

// WithErrorReason sets the ErrorReason field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ErrorReason field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithErrorReason(value v1beta1.MachineStatusError) *MachineStatusApplyConfiguration {
	b.ErrorReason = &value
	return b
}

// WithErrorMessage sets the ErrorMessage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ErrorMessage field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithErrorMessage(value string) *MachineStatusApplyConfiguration {
	b.ErrorMessage = &value
	return b
}

// WithProviderStatus sets the ProviderStatus field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ProviderStatus field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithProviderStatus(value runtime.RawExtension) *MachineStatusApplyConfiguration {
	b.ProviderStatus = &value
	return b
}

// WithAddresses adds the given value to the Addresses field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Addresses field.
func (b *MachineStatusApplyConfiguration) WithAddresses(values ...v1.NodeAddress) *MachineStatusApplyConfiguration {
	for i := range values {
		b.Addresses = append(b.Addresses, values[i])
	}
	return b
}

// WithLastOperation sets the LastOperation field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the LastOperation field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithLastOperation(value *LastOperationApplyConfiguration) *MachineStatusApplyConfiguration {
	b.LastOperation = value
	return b
}

// WithPhase sets the Phase field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Phase field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithPhase(value string) *MachineStatusApplyConfiguration {
	b.Phase = &value
	return b
}

// WithConditions adds the given value to the Conditions field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Conditions field.
func (b *MachineStatusApplyConfiguration) WithConditions(values ...*ConditionApplyConfiguration) *MachineStatusApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithConditions")
		}
		b.Conditions = append(b.Conditions, *values[i])
	}
	return b
}

// WithAuthoritativeAPI sets the AuthoritativeAPI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AuthoritativeAPI field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithAuthoritativeAPI(value v1beta1.MachineAuthority) *MachineStatusApplyConfiguration {
	b.AuthoritativeAPI = &value
	return b
}

// WithSynchronizedGeneration sets the SynchronizedGeneration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SynchronizedGeneration field is set to the value of the last call.
func (b *MachineStatusApplyConfiguration) WithSynchronizedGeneration(value int64) *MachineStatusApplyConfiguration {
	b.SynchronizedGeneration = &value
	return b
}