# Semantic Release Setup Guide

This project uses [semantic-release](https://github.com/semantic-release/semantic-release) for automated versioning and changelog generation.

## 📋 How It Works

Semantic-release automatically:
- Determines the next version number based on commit messages
- Generates a changelog from commit messages
- Creates a GitHub release with release notes
- Tags the release in Git

## 🚀 Setup Instructions

### 1. Install Dependencies

```bash
npm install
```

### 2. Commit Message Format

Use the [Angular Commit Message Convention](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#commit):

```
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

### Types and Version Bumps

| Type | Description | Version Bump | Example |
|------|-------------|--------------|---------|
| `feat` | New feature | **Minor** (0.x.0) | `feat: add user authentication` |
| `fix` | Bug fix | **Patch** (0.0.x) | `fix: resolve nil pointer error` |
| `perf` | Performance improvement | **Patch** (0.0.x) | `perf: optimize database queries` |
| `refactor` | Code refactoring | **Patch** (0.0.x) | `refactor: restructure service layer` |
| `docs` | Documentation changes | **Patch** (0.0.x) | `docs: update API documentation` |
| `chore` | Maintenance tasks | **No release** | `chore: update dependencies` |
| `style` | Code style changes | **No release** | `style: format code` |
| `test` | Adding/updating tests | **No release** | `test: add unit tests for service` |
| **BREAKING CHANGE** | Breaking API change | **Major** (x.0.0) | See below |

### Breaking Changes

Add `BREAKING CHANGE:` in the commit body or footer:

```bash
git commit -m "feat: change API response format

BREAKING CHANGE: API now returns data in a different structure"
```

Or use `!` after the type:

```bash
git commit -m "feat!: change authentication method"
```

## 📝 Commit Examples

### Feature (Minor Release)
```bash
git commit -m "feat(auth): add JWT authentication middleware"
```

### Bug Fix (Patch Release)
```bash
git commit -m "fix(db): handle connection timeout errors"
```

### Breaking Change (Major Release)
```bash
git commit -m "feat(api): redesign REST API structure

BREAKING CHANGE: All API endpoints have been restructured. 
See migration guide for details."
```

### Documentation (Patch Release)
```bash
git commit -m "docs: add setup instructions for PostgreSQL"
```

### Chore (No Release)
```bash
git commit -m "chore: update go dependencies"
```

## 🔄 Release Workflow

### Automatic Release (Recommended)

1. **Make changes** to your code
2. **Commit** with proper conventional commit format:
   ```bash
   git add .
   git commit -m "feat: add new feature"
   ```
3. **Push to main branch**:
   ```bash
   git push origin main
   ```
4. **GitHub Actions** will automatically:
   - Analyze commits
   - Determine version bump
   - Generate CHANGELOG.md
   - Create GitHub release
   - Tag the release

### Manual Release (Local)

If you want to test semantic-release locally:

```bash
# Dry run (no actual release)
npx semantic-release --dry-run

# Actual release (requires GITHUB_TOKEN)
export GITHUB_TOKEN=your_github_token
npx semantic-release
```

## 🔑 GitHub Token Setup

For the GitHub Action to work, you need:

1. The default `GITHUB_TOKEN` is automatically provided by GitHub Actions
2. Make sure the workflow has write permissions (already configured in `release.yml`)

### For Manual/Local Releases

1. Go to GitHub Settings → Developer settings → Personal access tokens
2. Generate a new token with these permissions:
   - `repo` (all)
   - `write:packages`
3. Set the token as environment variable:
   ```bash
   export GITHUB_TOKEN=your_token_here
   ```

## 📊 Version Strategy

Current configuration:
- **Major** (1.0.0 → 2.0.0): Breaking changes
- **Minor** (1.0.0 → 1.1.0): New features, refactors, performance improvements
- **Patch** (1.0.0 → 1.0.1): Bug fixes, documentation updates

## 🎯 Best Practices

1. **Write clear commit messages** - They become your changelog
2. **One logical change per commit** - Makes history cleaner
3. **Use scopes** - Help organize changes (e.g., `feat(auth):`, `fix(db):`)
4. **Test before pushing** - Broken code means broken releases
5. **Review CHANGELOG.md** - After release to ensure quality

## 🛠️ Configuration Files

- `.releaserc.json` - Semantic-release configuration
- `package.json` - Node.js dependencies
- `.github/workflows/release.yml` - GitHub Actions workflow

## 📖 Common Scopes

Suggested scopes for this project:
- `auth` - Authentication/authorization
- `api` - API endpoints
- `db` - Database operations
- `config` - Configuration
- `dto` - Data transfer objects
- `domain` - Domain services
- `repo` - Repositories
- `middleware` - Middleware
- `utils` - Utility functions
- `docs` - Documentation
- `ci` - CI/CD related

## 🐛 Troubleshooting

### Release not triggered
- Check commit message format
- Ensure you pushed to `main` branch
- Verify GitHub Actions has permissions
- Check for `[skip ci]` in commit message

### Wrong version bump
- Review commit message types
- Check for unintended `BREAKING CHANGE` markers
- Use `--dry-run` to preview

### CHANGELOG not updated
- Ensure `@semantic-release/changelog` is installed
- Check `.releaserc.json` configuration
- Verify git assets configuration

## 📚 Additional Resources

- [Semantic Release Documentation](https://semantic-release.gitbook.io/semantic-release/)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [Angular Commit Guidelines](https://github.com/angular/angular/blob/master/CONTRIBUTING.md#commit)

## 🔄 Migration from Manual Versioning

If migrating from manual versioning:

1. Ensure your current version in git tags matches your actual version
2. Follow conventional commits from now on
3. Let semantic-release handle all future releases
4. Update documentation to reference the new workflow

---

**Remember**: Good commit messages = Good changelogs = Happy users! 🎉
