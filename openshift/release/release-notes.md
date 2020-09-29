# Tekton Triggers v0.7.0

## Features

* Add CEL function to parse YAML ([#636](https://github.com/tektoncd/triggers/pull/636))

This PR will add a CEL function named parseYAML that can parse a YAML string into a map of strings to dynamic values
Syntax: .parseYAML() -> map<string, dyn>

* Improve the error from the context evaluation ([#646](https://github.com/tektoncd/triggers/pull/646))

Improve the error messages around parsing CEL expressions.
This improves the granularity of the error messages that are returned by evaluating expressions.
It also improves the error message when parsing the hook body when creating the evaluation environment.

* Support for marshaling other types in CEL ([#668](https://github.com/tektoncd/triggers/pull/668))

This adds support for marshaling bool values, and maps if they're used as the
values of expressions in a CEL overlay.

* Add nodeSelector and replicas to Eventlistener ([#625](https://github.com/tektoncd/triggers/pull/625))

Add nodeSelector and replicas feature to eventListener. With this, user could schedule eventListener pod to the node with
specific label. Also, if needed, user could specify the number of replicas in yaml file.

* Provide the incoming EventListener URL to the Webhook Interceptor ([#669](https://github.com/tektoncd/triggers/pull/669))

Webhook Interceptors can parse the EventListener-Request-URL if they want to
extract parameters from the original request URL being handled by the
EventListener.

## Breaking Changes :rotating_light:
* Remove deprecated $(params) ([#690](https://github.com/tektoncd/triggers/pull/690))

This is a breaking change as this PR remove complete support of $(params) and moved to $(tt.params) in order to avoid confusion between resourcetemplates and triggertemplate params

## Fixes :bug:
* Fix getting-started triggers ([#642](https://github.com/tektoncd/triggers/pull/642))
The EventListener was referring to the Binding via name instead of ref. Also, run the getting-started
examples as part of the e2e YAML tests. While this won't catch all issues with the examples, it should
catch obvious syntax issues like this one.

* Pass url through ([#657](https://github.com/tektoncd/triggers/pull/657))
Fix a bug in the sink where is not passing the URL through to the incoming requests.

* Fix triggertemplate validation to validate missing spec field ([#691](https://github.com/tektoncd/triggers/pull/691))

## Misc :hammer:
* Use sets.NewString instead of map[string]struct{} ([#663](https://github.com/tektoncd/triggers/pull/663))
* Update to pipeline knative 0.15 ([#661](https://github.com/tektoncd/triggers/pull/661))
* Update tektoncd/pipeline to v0.14.2 ([#684](https://github.com/tektoncd/triggers/pull/684))
* Update golang.org/x/text to v0.3.3 ([#674](https://github.com/tektoncd/triggers/pull/674))

## Docs :book:
* Add cel filter for pull request actions in github example ([#637](https://github.com/tektoncd/triggers/pull/637))
* Update docs and examples to use ref instead of name for bindings ([#645](https://github.com/tektoncd/triggers/pull/645))
* Add EventListener Response in the Doc ([#664](https://github.com/tektoncd/triggers/pull/664))
* Remove unused Ref from EventListenerTrigger ([#677](https://github.com/tektoncd/triggers/pull/677))

## How to upgrade from v0.6.1 :up_arrow:
1. Change any $(params) to $(tt.params) in TriggerTemplate
2. Install Triggers. One liner:
```text
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/previous/v0.7.0/release.yaml
```

# Tekton Triggers v0.8.0

## Features
* Propagate annotations from Eventlistener to service and deployment ([#712](https://github.com/tektoncd/triggers/pull/712))

This PR adds feasibility to propagate annotations from the EventListener to deployment and services.
If there are any custom annotations on services/deployment, then it needs to be added to EventListener annotations, so that those will be propagated otherwise they will be overwritten.

* Add validation for replicas ([#717](https://github.com/tektoncd/triggers/pull/717))

This PR handles proper validation for replicas which are provided as part of EventListener spec.
Produce below error if provided replica is invalid
```text
Error from server (BadRequest): error when creating "STDIN": admission webhook "validation.webhook.triggers.tekton.dev" 
```

* Add TriggerCRD object Ref to Eventlistener Spec ([#726](https://github.com/tektoncd/triggers/pull/726))

This PR helps to specify TriggerCRD object inside the Eventlistener spec as a reference using triggerRef field,
So this way user can create TriggerCRD separately and bind it inside Eventlistener spec.

* Add validation and default for TriggerCRD object ([#738](https://github.com/tektoncd/triggers/pull/738))

This PR adds the validation and defaults around TriggerCRD.

## Breaking Changes :rotating_light:
* Switch trigger sa based auth to impersonate ([#705](https://github.com/tektoncd/triggers/pull/705))
* Switch trigger sa ref from global to namespace scoped ([#704](https://github.com/tektoncd/triggers/pull/704))

The optional EventListenerTrigger based level of authentication for creating Tekton object has had its ServiceAccount reference changed from an ObjectReference to a string ServiceAccountName, effectively enforcing that the ServiceAccount be in the same namespace as the EventListenerTrigger.

## Fixes :bug:
* Fix update deployment when there is a change in replicas ([#715](https://github.com/tektoncd/triggers/pull/715))
* Add basic syntactical parsing of CEL filters and expressions ([#745](https://github.com/tektoncd/triggers/pull/745))

Perform simple syntax checking of CEL filter and overlays in the Webhook validator, perfunctory syntax validation of the expressions in the interceptor but it won't detect logical errors (expressions that rely on JSON bodies).

## Docs :book:
* Update eventlistener doc to include eventlistener responsibility ([#721](https://github.com/tektoncd/triggers/pull/721))

Added information about the responsibility of Eventlistener.

* Update triggertemplate doc and example to use tt.params ([#725](https://github.com/tektoncd/triggers/pull/725))
* Clarify Bitbucket and Github HMAC generation instruction in example ([#728](https://github.com/tektoncd/triggers/pull/728))
* Update pull request template for release-note ([#743](https://github.com/tektoncd/triggers/pull/743))

## How to upgrade from v0.7.0 :up_arrow:
```text
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/previous/v0.8.0/release.yaml
```
`NOTE: Due to #750, you will have to manually add replicas: 1 to your EL spec before applying the upgrade`

# Tekton Triggers v0.8.1

## Fixes :bug:
* Merge annotations before propagation ([#753](https://github.com/tektoncd/triggers/pull/753))

Triggers no longer overwrites annotations set on the underlying deployment and service objects.

* Add a Idle Timeout for EventListener sink ([#755](https://github.com/tektoncd/triggers/pull/755))

EventListeners close idle connections after 120 seconds.

## How to upgrade from v0.8.0 :up_arrow:
```text
kubectl apply -f https://storage.googleapis.com/tekton-releases/triggers/previous/v0.8.1/release.yaml
```