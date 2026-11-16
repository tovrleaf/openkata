# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.0.0] - 2026-11-16

### Added

- Check script that flags long, tangled and machine-sounding
  sentences in Markdown, text and templ files, for files, folders,
  pasted text or a branch's changed lines
- Watch list of inflated words and phrases with plain alternatives
- Workflow to judge each flagged sentence and rewrite only those in
  place
- Optional git hooks that check the staged text and the commit
  message on every commit
