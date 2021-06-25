# Tekton Triggers v0.14.1

## Deprecation Notices

* :rotating_light: Deprecate podTemplate field (#1102)
Handled deprecation of `podTemplate` properly so that there won't be any issue with Upgrade which we faced as part of `0.14.0`.
Reference Issue https://github.com/tektoncd/triggers/issues/1098

# Tekton Triggers v0.14.0

## Features

* :sparkles: Adding opencensus metrics (#1061)

Added Eventlistener OpenCensus metrics which captures metrics at process level.

## Deprecation Notices

* :rotating_light: Add a Ready StatusCondition for EventListener (#1082)

The EventListener Status now has a new Condition called Ready. Ready is set to True when the other status conditions are also true (i.e. the EventListener is ready to serve traffic). It is false if any of the other statuses are false.

**DEPRECATION NOTICE:** In the future, we plan to deprecate the other status conditions in favor of the Ready status condition.

* :rotating_light: Add EventListener UID to sink response, mark name/namespace as deprecated (#1087)

`eventListener` and `namespace` fields in EventListener response are now deprecated. Use `eventListenerUID` instead.

## Backwards incompatible changes

* :rotating_light: Removed deprecated fields (#1040)

As part of this release we have removed deprecated fields `ServiceType` and `PodTemplate` from the EventListener Spec.

# Tekton Triggers v0.13.0

## Features

* :sparkles: Adding label selector for triggers [#970](https://github.com/tektoncd/triggers/pull/970)

    Users can now configure triggers for an EventListener using labels

*  :sparkles:  Add ClusterInterceptor CRD for registering interceptors (#960)

    A new CRD called ClusterInterceptor has been added that allows for users to register new pluggable Interceptor types

*  :sparkles:  Add API to configure Interceptors from a Trigger (#1001)
    
    Trigger spec authors can now configure interceptors using a new API that includes  a `ref` field to refer to a ClusterInterceptor, and 
    a `params` field to add parameters to pass on to the interceptor for processing

*  :sparkles:  Migrate core interceptors to use InterceptorType CRD (#976)

    The four bundled interceptors (CEL, GitHub, GitLab, BitBucket) are now implemented using the new ClusterInterceptor CRD

*  :sparkles:   Migrate core interceptors to new format (#1029)

     Any new Triggers created using the old style syntax for core interceptors is now automatically switched to the new refs/params 
     based syntax.


## Deprecation Notices

* :rotating_light: Move replicas to KubernetesResource from eventlistener spec (#1021)

    Deprecated replicas from EventListener spec level and added replicas to KubernetesResource.

* :rotating_light:  Migrate core interceptors to new format (#1029)

    Interceptors in a Trigger are now configured using a new `ref` and `params` based syntax as described in #1001, and in #976 we 
    have implemented the core interceptors using the new ClusterInterceptor CRD. So, these core interceptors can now be configured 
    using the new syntax. The old way of configuring these interceptors is now deprecated but it will continue to work until it is removed 
    in a future release. The defaulting webhook will automatically switch any usages of the old syntax to the new for any new Triggers.
