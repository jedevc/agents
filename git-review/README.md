# Git review module

This Dagger module provides automated code review functionality for Git pull requests.

Features:

- Automated review of Git pull request changes.
- Diff analysis between base and target branches.
- AI-powered review comments and recommendations.
- Support for custom review instructions.

Example usage to review a PR:

    dagger call review --base https://github.com/dagger/dagger.git --ref https://github.com/dagger/dagger.git#refs/pull/<pr>/head

