# Tekton Triggers v0.12.0

## Features 

* Add support for custom object to triggers eventlistener ([#958](https://github.com/tektoncd/triggers/pull/958))
    
    Introduced new field customResource to support Knative Service for EventListener
    
    ```
    apiVersion: triggers.tekton.dev/v1alpha1
    kind: EventListener
    metadata:
      name: github-listener-interceptor-customresource
    spec:
      ...
      resources:
        customResource:
          apiVersion: serving.knative.dev/v1
          kind: Service
          metadata:
          spec:
            template:
              spec:
                serviceAccountName: tekton-triggers-example-sa
                containers:
                - resources:
                    requests:
                      memory: "64Mi"
                      cpu: "250m"
                    limits:
                      memory: "128Mi"
                      cpu: "500m"
    ```
        
* Validate Event Body for Json Format ([#969](https://github.com/tektoncd/triggers/pull/969))

    We now throw http.BadRequest status code(400) if event payload isn't json.

## Backwards incompatible changes :rotating_light:

In current release:

* Remove deprecated field template.Name in favour of template.Ref ([#919](https://github.com/tektoncd/triggers/pull/919))

    Deprecated field template.Name in has been removed in favor of template.Ref

* Switch to UUID for event IDs ([#926](https://github.com/tektoncd/triggers/pull/926))

    Change the event ID representation from a 5 character random string to a UUID. 

# Tekton Triggers v0.11.0

## Features :sparkles:
* Migrate GitLab, BitBucket, GitHub interceptors to new interface ([#832](https://github.com/tektoncd/triggers/pull/832))
* Implement marshalJSON CEL function ([#842](https://github.com/tektoncd/triggers/pull/842))
  
  New CEL function `marshalJSON` that can encode a JSON object or array to a string.
* Add a server for serving core interceptors ([#858](https://github.com/tektoncd/triggers/pull/858))
  
  Add a HTTP handler for serving core interceptors and this packages all 4 core interceptors into a single HTTP server.
  Each interceptor is available at a different path e.g. /cel for CEL etc.
* Move core interceptors to their own server ([#878](https://github.com/tektoncd/triggers/pull/878))  
## Deprecation Notices :rotating_light:
* Deprecate PodTemplate and ServiceType in favour of Resource ([#897](https://github.com/tektoncd/triggers/pull/897))

  Deprecate PodTemplate and ServiceType in favour of Resource.
## Backwards incompatible changes :rotating_light:
* Remove the template.Name field ([#898]((https://github.com/tektoncd/triggers/pull/898)))

Action required: The Template.Name field has been removed from the Trigger Spec. Please use Template.Ref instead.
* Remove deprecated spec style embedded bindings ([#900](https://github.com/tektoncd/triggers/pull/900))

BREAKING CHANGE:
The `Spec` field has been removed from the `TriggerSpecBinding`.
