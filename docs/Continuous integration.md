# Continuous integration

## Complete CI

```mermaid
mindmap
(fullproxy)
	[main]
		Push
			Release
	[dev]
		Pull Request
			Test
		Push
			Tagging
```

## Pipes

### Test

| Feature     | Value                             |
| ----------- | --------------------------------- |
| Executes    | On **PR** to `dev`                |
| Permissions | **Read only** repository contents |

```mermaid
stateDiagram-v2
	direction LR
	state Build {
		direction LR
        [*] --> BuildBinary
        BuildBinary --> [*]
    }
    state Test {
		direction LR
        [*] --> DockerCompose
        DockerCompose --> UnitTest
        UnitTest --> [*]
    }
	[*] --> Build
	Build --> Test
	Test --> [*]
```

### Tagging

| Feature     | Value                              |
| ----------- | ---------------------------------- |
| Executes    | On **Pushes** to `dev`             |
| Permissions | **Read/Write** repository contents |

```mermaid
stateDiagram-v2
	direction LR
	
	state Tagging {
		direction LR
        [*] --> CheckoutDev
        CheckoutDev --> StandardVersion
        StandardVersion --> Push
        Push --> [*]
    }
    
    [*] --> Tagging
    Tagging --> [*]
```

### Release

| Feature     | Value                                                        |
| ----------- | ------------------------------------------------------------ |
| Executes    | On **Pushes** to `main`                                      |
| Permissions | **Read only** repository contents, **Read/Write** releases, **Read/Write** packages |

```mermaid
stateDiagram-v2
	direction LR
	
	state Versioning {
		direction LR
		VersionPy: Read tag version
		[*] --> VersionPy
		VersionPy --> [*]
	}
	
	state Build {
		direction LR
		state fork_build <<fork>>
		state join_build <<join>>
		[*] --> fork_build
		fork_build --> Windows
		fork_build --> Linux
		fork_build --> FreeBSD
		fork_build --> OSX
		Windows --> join_build
		Linux --> join_build
		FreeBSD --> join_build
		OSX --> join_build
		join_build --> [*]
	}
	
	state Release {
		direction LR
		state fork_release <<fork>>
		state join_release <<join>>
		GitHub: GitHub Releases
		Pages: GitHub pages
		[*] --> fork_release
		fork_release --> GitHub
		fork_release --> Pages
		GitHub --> join_release
		Pages --> join_release
		join_release --> [*]
	}
	
	state UPX {
		direction LR
		upx: Run UPX
		[*] --> upx
		upx --> [*]
	}
	
	[*] --> Versioning
	Versioning --> Build
	Build --> UPX
	UPX --> Release
	Release --> [*]
```

