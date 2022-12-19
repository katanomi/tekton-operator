/*
Copyright 2021 The Tekton Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"context"

	"knative.dev/pkg/ptr"
)

const (
	DefaultMetricsPipelinerunLevel       = "pipeline"
	DefaultMetricsTaskrunLevel           = "task"
	DefaultMetricsPipelierunDurationType = "histogram"
	DefaultMetricsTaskrunDurationType    = "histogram"
	KatanomiMigratePipelineApiFields     = "katanomi.dev/enable-api-fields-migrate"
)

func (tp *TektonPipeline) SetDefaults(ctx context.Context) {
	tp.Spec.PipelineProperties.setDefaults()
	tp.EnableApiFields(ctx)
}

func (tp *TektonPipeline) EnableApiFields(ctx context.Context) {
	// as default, enable-api-fields would be alpha
	if tp.Spec.EnableApiFields == "" {
		tp.Spec.EnableApiFields = ApiFieldAlpha
	}

	// for old tektonpipeline resource, we will force change enable-api-fields to alpha
	if tp.Annotations == nil {
		tp.Annotations = map[string]string{}
	}
	if _, ok := tp.Annotations[KatanomiMigratePipelineApiFields]; !ok {
		// we change EnableApiFields to alpha,
		// avoid upgrade from old version that cannot change to alpha from stable
		tp.Spec.PipelineProperties.EnableApiFields = ApiFieldAlpha
		tp.Annotations[KatanomiMigratePipelineApiFields] = "true"
	}

}

func (p *PipelineProperties) setDefaults() {
	// Disabling this for now and will be removed in next release
	// disabling will hide this from users in TektonConfig/TektonPipeline
	if p.DisableHomeEnvOverwrite != nil {
		p.DisableHomeEnvOverwrite = nil
	}
	if p.DisableWorkingDirectoryOverwrite != nil {
		p.DisableWorkingDirectoryOverwrite = nil
	}

	if p.DisableCredsInit == nil {
		p.DisableCredsInit = ptr.Bool(false)
	}
	if p.RunningInEnvironmentWithInjectedSidecars == nil {
		p.RunningInEnvironmentWithInjectedSidecars = ptr.Bool(true)
	}
	if p.RequireGitSshSecretKnownHosts == nil {
		p.RequireGitSshSecretKnownHosts = ptr.Bool(false)
	}
	if p.EnableTektonOciBundles == nil {
		p.EnableTektonOciBundles = ptr.Bool(false)
	}
	if p.EnableCustomTasks == nil {
		p.EnableCustomTasks = ptr.Bool(false)
	}

	if p.ScopeWhenExpressionsToTask == nil {
		p.ScopeWhenExpressionsToTask = ptr.Bool(false)
	}
	if p.MetricsPipelinerunDurationType == "" {
		p.MetricsPipelinerunDurationType = DefaultMetricsPipelierunDurationType
	}
	if p.MetricsPipelinerunLevel == "" {
		p.MetricsPipelinerunLevel = DefaultMetricsPipelinerunLevel
	}
	if p.MetricsTaskrunDurationType == "" {
		p.MetricsTaskrunDurationType = DefaultMetricsTaskrunDurationType
	}
	if p.MetricsTaskrunLevel == "" {
		p.MetricsTaskrunLevel = DefaultMetricsTaskrunLevel
	}

}
