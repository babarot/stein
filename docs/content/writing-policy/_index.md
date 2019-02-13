---
title: "Writing Policy"
date: 2017-10-17T15:26:15Z
draft: false
weight: 30

---

This section covers how to write Stein policies.

It is recommended that you complete the [getting started guide]({{< relref "getting-started" >}}) prior to reading this section of the documentation. The getting started guide will explain all the basics of writing and testing policies. This section of the documentation will serve as more of a reference guide to all available features of Stein.

Stein provides a language and workflow for building policy across any system that embeds Stein. By learning Stein once, you are able to effectively control access to many systems using Stein's powerful features. Stein also provides a [local CLI]() called the Stein Simulator for developing and testing Stein policies. This CLI can be integrated into a daily developer's workflow along with CI to treat [policy as code](policy-as-code.md).

Stein uses its own [language]() for writing policies. You can view a [language reference]() as well as the [specification]() for details. You don't have to read those documents immediately, since the language should be easy enough to pick up throughout this section.

## What's next

- [Syntax](syntax/syntax.md)
  - [Interpolation Syntax](syntax/interpolation.md)
- [Policy](policy.md)
  - [Rules](policy/rules.md)
  - [Variables](policy/variables.md)
  - [Functions](policy/functions.md)
  - [Config](policy/config.md)
