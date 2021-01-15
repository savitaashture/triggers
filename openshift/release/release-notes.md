# Tekton Triggers v0.10.2

## Fixes :bug:
* Merge extensions into body for webhook interceptors ([#860](https://github.com/tektoncd/triggers/pull/860))

  Extensions added by a CEL Interceptor will be passed on to webhook interceptors by merging the extension fields into the event body under a `extensions` field. 
  See docs/eventlisteners.md ##chaining-interceptors for more details.

# Tekton Triggers v0.10.0

## Features :sparkles:

* Allow users to set resources as part of podtemplate ([#815](https://github.com/tektoncd/triggers/pull/815))

  Now trigger allow users to specify their resource information in eventlistener
  
```yaml
apiVersion: triggers.tekton.dev/v1alpha1
kind: EventListener
metadata:
  name: github-listener-interceptor
spec:
  ...
  resources:
    kubernetesResource:
      spec:
        template:
          spec:
            serviceAccountName: tekton-triggers-github-sa
            containers:
              - resources:
                  requests:
                    memory: "64Mi"
                    cpu: "250m"
                  limits:
                    memory: "128Mi"
                    cpu: "500m"
```  
* Use Listers to fetch data in the Sink([#821](https://github.com/tektoncd/triggers/pull/821))
  
  EventListener ServiceAccounts now need to have "list" and "watch" verbs in addition to "get" for all triggers resources. See examples at https://github.com/tektoncd/triggers/tree/v0.10.0/examples/role-resources/triggerbinding-roles/role.yaml and https://github.com/tektoncd/triggers/tree/v0.10.0/examples/role-resources/clustertriggerbinding-roles/clusterrole.yaml

* TEP-0022: Switch to immutable input event bodies ([#828](https://github.com/tektoncd/triggers/pull/828))
  Migrate CEL to new Interceptor Interface

* Add EventListener Selector For TriggerCRD ([#773](https://github.com/tektoncd/triggers/pull/773))

  Added Namespace Selector field for EventListener which enables EventListener to serve across the namespace. namespaceSelector field with matchNames need to be provided to enable selector.
  ```yaml
  namespaceSelector:
    matchNames:
    - nsName1
    - nsName2
  ```    
  If namespace selector is used, the service account for the EventListener will need a clusterRole. See the example at https://github.com/tektoncd/triggers/tree/v0.10.0/examples/selectors/01_rbac.yaml

* Allow secure connection to eventlistener pod ([#819](https://github.com/tektoncd/triggers/pull/819))

  HTTPS connection to eventlistener can be configured by tweaking eventlistener configuration
  ```yaml
   apiVersion: triggers.tekton.dev/v1alpha1
   kind: EventListener
   metadata:
     name: github-listener-interceptor
   spec:
     ...
     resources:
       kubernetesResource:
         spec:
           template:
             spec:
               serviceAccountName: tekton-triggers-github-sa
               containers:
               - env:
                 - name: TLS_CERT
                   valueFrom:
                     secretKeyRef:
                       name: tls-key-secret
                       key: tls.crt
                 - name: TLS_KEY
                   valueFrom:
                     secretKeyRef:
                       name: tls-key-secret
                       key: tls.key
  ```
* Drop escaping of strings in the JSON ([#823](https://github.com/tektoncd/triggers/pull/823))
  Change the escaping of parameters into TriggerTemplates.

## Backwards incompatible changes :rotating_light:
In the current release:

* TEP-0022: Switch to immutable input event bodies ([#828](https://github.com/tektoncd/triggers/pull/828))

  action required: If you are using overlays in the CEL Interceptor, please update your bindings to use $(extensions.) instead of $(body.)

  BREAKING CHANGE:
  
  CEL overlays now add fields to a new top level extensions field instead of the modifying the incoming event body. TriggerBindings can access values within this new extensions field using `$(extensions.<key>)` syntax.

* Drop escaping of strings in the JSON ([#823](https://github.com/tektoncd/triggers/pull/823))

  Previously, parameters were escaped as they were being replaced into a TriggerTemplate, by replacing double-quotes " with an escaped version ", this functionality has been removed, as it was breaking quoted strings, and in some cases, rendering the resulting output unparseable.

  action required: If you were relying on the escaping, you can retain the old behaviour by adding an annotation to an affected TriggerTemplate, `triggers.tekton.dev/old-escape-quotes: "true"`

# Tekton Triggers v0.9.0

## Features :sparkles:

* Enhance existing eventlistener to support PodTemplate for Deployment using duck type ([#734](https://github.com/tektoncd/triggers/pull/734))

  A new field `resources` has been introduced as part of EventListener spec
```yaml
apiVersion: triggers.tekton.dev/v1alpha1
kind: EventListener
metadata:
  name: github-listener-interceptor
spec:
  triggers:
    ...
  resources:
    kubernetesResource:
      serviceType: NodePort
      spec:
        template:
          metadata:
            labels:
              key: "value"
            annotations:
              key: "value"
          spec:
            serviceAccountName: tekton-triggers-github-sa
            nodeSelector:
              app: test
            tolerations:
              - key: key
                value: value
                operator: Equal
                effect: NoSchedule
```
As of now the `resources` field supports `kubernetesResource` which helps us use PodSpecable ducktype. 
For backward compatibility both ways are supported and `resources` are optional.

* cel-go with string upper and lower-casing. ([#766](https://github.com/tektoncd/triggers/pull/766))
  
  Updated version of cel-go with new functionality for upper and lower-casing ASCII strings. e.g. body.upperMsg.lowerAscii()

* Adds support for name/value embedded Bindings ([#768](https://github.com/tektoncd/triggers/pull/768))
  
  TriggerBindings can now be embedded by using just name/value fields inside a Trigger or a EventListener.
  
  BREAKING CHANGE: With this change, users cannot specify both name and ref
  for a single binding. Use `ref` to refer to a TriggerBinding resource and
  `name` for embedded bindings.
```yaml
# NEW SYNTAX:
bindings:
- ref: some-name
- name: commit_id # embedded binding
  value: "$(body.head_commit_id)"

# OLD SYNTAX:
bindings:
- name: some-name
  spec:
    params:
    - name: commit_id # embedded binding
      value: "$(body.head_commit_id)
```

* PodSecurityPolicy Config fixes to allow running in restricted envs ([#707](https://github.com/tektoncd/triggers/pull/707))
  
  The PodSecurityPolicy configuration was updated so that container must run as non-root, and the RBAC for utilizing the PSP was moved from cluster scoped to namespace scoped to better ensure triggers can not utilize other PSPs unrelated to this project

* Add support for embedded trigger templates ([#783](https://github.com/tektoncd/triggers/pull/783))
  
  Users can now specify embed a TriggerTemplate spec inside a Trigger.
  
## Deprecation Notices :rotating_light:
* Use template.ref instead of template.name ([#787](https://github.com/tektoncd/triggers/pull/787))

  The template.name field is deprecated (in favor of template.ref) and will be removed in a future release.

* Adds support for name/value embedded Bindings ([#768](https://github.com/tektoncd/triggers/pull/768))

  The old syntax for embedded TriggerBindings (using `spec.params`) is deprecated. Use the new name/value fields instead:  
  
## Backwards incompatible changes :rotating_light:
In current release:

* Adds support for name/value embedded Bindings ([#768](https://github.com/tektoncd/triggers/pull/768))

  BREAKING CHANGE: users cannot specify both `name` and `ref` for a single binding. Use `ref` to refer to a TriggerBinding resource and `name` for embedded bindings.

* remove interceptor cross namespace secret references ([#748](https://github.com/tektoncd/triggers/pull/748))

  BREAKING CHANGE: any interceptors that attempt to reference a Secret outside of the EventListener's namespace will have to change to have those Secrets reside in the EventListener's namespace.

* Apply triggertemplate param default value if triggerbinding param missing from body/header ([#761](https://github.com/tektoncd/triggers/pull/761))

  Triggers now applies the default value of a trigger template param if it's value cannot be resolved from a TriggerBinding. 
  This is a BREAKING CHANGE. Previously, Triggers would throw an error saying that the binding value could not be resolved (See #568)
