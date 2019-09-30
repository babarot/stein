---
title: "Policy as Code"
date: 2019-01-17T15:26:15Z
draft: false
weight: 12

---

Policy as code is the idea of writing code in a high-level language to manage and automate policies. By representing policies as code in text files, proven software development best practices can be adopted such as version control, automated testing, and automated deployment.

Many existing policy or ACL systems do not practice policy as code. Many policies are set by clicking in a GUI, which isn't easily repeatable nor versionable. They usually don't provide any system for testing policies other than testing an action that would violate the policy. This makes it difficult for automated testing. And the policy language itself varies by product.

Stein is built around the idea and provides all the benefits of policy as code.

!!! Note
    The idea of "Policy as Code" is proposed by HashiCorp and HashiCorp Sentinel. [Policy as Code - Sentinel by HashiCorp](https://docs.hashicorp.com/sentinel/concepts/policy-as-code)
