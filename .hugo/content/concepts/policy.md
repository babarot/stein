---
title: "Policy"
date: 2019-01-17T15:26:15Z
draft: false
weight: 11

---

## What's a policy?

The policy describes the rule "The configuration file should be written like this". How the configuration file should be written depends on the company and team. Stein respects and embraces the approaches adopted by Terraform and [HashiCorp Sentinel](https://docs.hashicorp.com/sentinel/intro/) so that it can be defined flexibly.

Stein can test against the configuration file based on those policies. Therefore, by combining with CI etc., you can enforce rules on the configuration file.

## Why needs policies?

Nowadays, thanks to the infiltration of the concept of "Infrastructure as Code", many infrastructure settings are coded in a configuration file language such as YAML.

For YAML files that we have to maintain like Kubernetes manifest files, as we continue to maintain them in the meantime, we will want to unify the writing style with policies like a style guide.

Let's say reviewing Kubernetes YAML. In many cases, you will find points to point out repeatedly in repeated reviews. For example, explicitly specifying namespace, label, and so on. The things that humans point out every time should be checked mechanically beforehand. It's important for reviewers as well as for reviewees.

Besides, there are points that can not be pointed out in the review, and a lot of difficult points to comment. For example, it's whether the specified namespace name is correct or not. In addition, let's say Terraform use-case.
As an example: before infrastructure as code and autoscaling, if an order came through for 5,000 new machines, a human would likely respond to the ticket verifying that the user really intended to order 5,000 new machines. Today, automation can almost always freely order 5,000 new compute instances without any hesitation, which can result in unintended expense or system instability ([HashiCorp Sentinel](https://docs.hashicorp.com/sentinel/intro/why) has basically the same intention and Stein's one also comes from that).

In order to avoid these accidents in advance, it is very important to define policies as codes and warn them based on them.
