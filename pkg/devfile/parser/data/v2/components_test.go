package v2

import (
	"reflect"
	"testing"

	v1 "github.com/devfile/api/pkg/apis/workspaces/v1alpha2"
	"github.com/devfile/library/pkg/testingutil"
)

func TestDevfile200_AddComponent(t *testing.T) {

	tests := []struct {
		name              string
		currentComponents []v1.Component
		newComponents     []v1.Component
		wantErr           bool
	}{
		{
			name: "case 1: successfully add the component",
			currentComponents: []v1.Component{
				{
					Name: "component1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "component2",
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			newComponents: []v1.Component{
				{
					Name: "component2",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "component3",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "case 2: error out on duplicate component",
			currentComponents: []v1.Component{
				{
					Name: "component1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
				{
					Name: "component2",
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			newComponents: []v1.Component{
				{
					Name: "component1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.currentComponents,
						},
					},
				},
			}

			got := d.AddComponents(tt.newComponents)

			if !tt.wantErr && got != nil {
				t.Errorf("TestDevfile200_AddComponents() unexpected error - %+v", got)
			} else if tt.wantErr && got == nil {
				t.Errorf("TestDevfile200_AddComponents() expected error but got nil")
			}

		})
	}
}

func TestDevfile200_UpdateComponent(t *testing.T) {

	tests := []struct {
		name              string
		currentComponents []v1.Component
		newComponent      v1.Component
	}{
		{
			name: "case 1: successfully update the component",
			currentComponents: []v1.Component{
				{
					Name: "Component1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image1",
							},
						},
					},
				},
				{
					Name: "component2",
					ComponentUnion: v1.ComponentUnion{
						Volume: &v1.VolumeComponent{},
					},
				},
			},
			newComponent: v1.Component{
				Name: "Component1",
				ComponentUnion: v1.ComponentUnion{
					Container: &v1.ContainerComponent{
						Container: v1.Container{
							Image: "image2",
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.currentComponents,
						},
					},
				},
			}

			d.UpdateComponent(tt.newComponent)

			components := d.GetComponents()

			matched := false
			for _, component := range components {
				if reflect.DeepEqual(component, tt.newComponent) {
					matched = true
					break
				}
			}

			if !matched {
				t.Error("TestDevfile200_UpdateComponent() error updating the component")
			}
		})
	}
}

func TestGetDevfileContainerComponents(t *testing.T) {

	tests := []struct {
		name                 string
		component            []v1.Component
		expectedMatchesCount int
	}{
		{
			name:                 "Case 1: Invalid devfile",
			component:            []v1.Component{},
			expectedMatchesCount: 0,
		},
		{
			name: "Case 2: Valid devfile with wrong component type (Openshift)",
			component: []v1.Component{
				{
					ComponentUnion: v1.ComponentUnion{
						Openshift: &v1.OpenshiftComponent{},
					},
				},
			},
			expectedMatchesCount: 0,
		},
		{
			name: "Case 3 : Valid devfile with correct component type (Container)",
			component: []v1.Component{
				testingutil.GetFakeContainerComponent("comp1"),
				testingutil.GetFakeContainerComponent("comp2"),
			},
			expectedMatchesCount: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.component,
						},
					},
				},
			}

			devfileComponents := d.GetDevfileContainerComponents()

			if len(devfileComponents) != tt.expectedMatchesCount {
				t.Errorf("TestGetDevfileContainerComponents error: wrong number of components matched: expected %v, actual %v", tt.expectedMatchesCount, len(devfileComponents))
			}
		})
	}

}

func TestGetDevfileVolumeComponents(t *testing.T) {

	tests := []struct {
		name                 string
		component            []v1.Component
		expectedMatchesCount int
	}{
		{
			name:                 "Case 1: Invalid devfile",
			component:            []v1.Component{},
			expectedMatchesCount: 0,
		},
		{
			name: "Case 2: Valid devfile with wrong component type (Kubernetes)",
			component: []v1.Component{
				{
					ComponentUnion: v1.ComponentUnion{
						Kubernetes: &v1.KubernetesComponent{},
					},
				},
			},
			expectedMatchesCount: 0,
		},
		{
			name: "Case 3: Valid devfile with correct component type (Volume)",
			component: []v1.Component{
				testingutil.GetFakeContainerComponent("comp1"),
				testingutil.GetFakeVolumeComponent("myvol", "4Gi"),
			},
			expectedMatchesCount: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.component,
						},
					},
				},
			}
			devfileComponents := d.GetDevfileVolumeComponents()

			if len(devfileComponents) != tt.expectedMatchesCount {
				t.Errorf("TestGetDevfileVolumeComponents error: wrong number of components matched: expected %v, actual %v", tt.expectedMatchesCount, len(devfileComponents))
			}
		})
	}

}

func TestGetPortExposure(t *testing.T) {
	urlName := "testurl"
	urlName2 := "testurl2"
	tests := []struct {
		name                string
		containerComponents []v1.Component
		wantMap             map[int]v1.EndpointExposure
		wantErr             bool
	}{
		{
			name: "Case 1: devfile has single container with single endpoint",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.PublicEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.PublicEndpointExposure,
								},
							},
						},
					},
				},
			},
		},
		{
			name:    "Case 2: devfile no endpoints",
			wantMap: map[int]v1.EndpointExposure{},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
						},
					},
				},
			},
		},
		{
			name: "Case 3: devfile has multiple endpoints with same port, 1 public and 1 internal, should assign public",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.PublicEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.PublicEndpointExposure,
								},
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.InternalEndpointExposure,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Case 4: devfile has multiple endpoints with same port, 1 public and 1 none, should assign public",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.PublicEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.PublicEndpointExposure,
								},
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.NoneEndpointExposure,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Case 5: devfile has multiple endpoints with same port, 1 internal and 1 none, should assign internal",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.InternalEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.InternalEndpointExposure,
								},
								{
									Name:       urlName,
									TargetPort: 8080,
									Exposure:   v1.NoneEndpointExposure,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Case 6: devfile has multiple endpoints with different port",
			wantMap: map[int]v1.EndpointExposure{
				8080: v1.PublicEndpointExposure,
				9090: v1.InternalEndpointExposure,
				3000: v1.NoneEndpointExposure,
			},
			containerComponents: []v1.Component{
				{
					Name: "testcontainer1",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName,
									TargetPort: 8080,
								},
								{
									Name:       urlName,
									TargetPort: 3000,
									Exposure:   v1.NoneEndpointExposure,
								},
							},
						},
					},
				},
				{
					Name: "testcontainer2",
					ComponentUnion: v1.ComponentUnion{
						Container: &v1.ContainerComponent{
							Container: v1.Container{
								Image: "image",
							},
							Endpoints: []v1.Endpoint{
								{
									Name:       urlName2,
									TargetPort: 9090,
									Secure:     true,
									Path:       "/testpath",
									Exposure:   v1.InternalEndpointExposure,
									Protocol:   v1.HTTPSEndpointProtocol,
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DevfileV2{
				v1.Devfile{
					DevWorkspaceTemplateSpec: v1.DevWorkspaceTemplateSpec{
						DevWorkspaceTemplateSpecContent: v1.DevWorkspaceTemplateSpecContent{
							Components: tt.containerComponents,
						},
					},
				},
			}
			mapCreated := d.GetPortExposure()
			if !reflect.DeepEqual(mapCreated, tt.wantMap) {
				t.Errorf("Expected: %v, got %v", tt.wantMap, mapCreated)
			}

		})
	}

}
