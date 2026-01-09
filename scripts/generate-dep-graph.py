#!/usr/bin/env python3
"""
Generate dependency graph for meme9 project.
Outputs Mermaid diagram, DOT file, and text representation.
"""

import os
import sys

# Service dependencies (service -> [dependencies])
SERVICE_DEPS = {
    'front8': ['auth-service', 'users-service', 'subscriptions-service', 'likes-service', 'posts-service', 'photos'],
    'posts-service': ['users-service', 'subscriptions-service', 'likes-service', 'auth-service'],
    'likes-service': ['users-service', 'auth-service'],
    'users-service': ['auth-service'],
    'subscriptions-service': ['auth-service'],
    'photos': ['auth-service'],
    'auth-service': [],
}

# External dependencies
EXTERNAL_DEPS = {
    'auth-service': ['MongoDB'],
    'users-service': ['MongoDB'],
    'subscriptions-service': ['MongoDB'],
    'likes-service': ['MongoDB'],
    'posts-service': ['MongoDB'],
    'photos': ['AWS S3'],
}

# Shared dependencies
SHARED_DEPS = {
    'api': ['auth-service', 'users-service', 'subscriptions-service', 'likes-service', 'posts-service', 'photos'],
}

def generate_mermaid():
    """Generate Mermaid diagram."""
    lines = ["graph TD"]
    
    # External services
    lines.append("    subgraph External[External Services]")
    lines.append("        MongoDB[MongoDB]")
    lines.append("        S3[AWS S3]")
    lines.append("    end")
    
    # Shared modules
    lines.append("    subgraph Shared[Shared Modules]")
    lines.append("        API[api<br/>Protobuf]")
    lines.append("    end")
    
    # Backend services
    lines.append("    subgraph Backend[Backend Services]")
    for service in ['auth-service', 'users-service', 'subscriptions-service', 'likes-service', 'posts-service', 'photos']:
        lines.append(f"        {service.replace('-', '_')}[{service}]")
    lines.append("    end")
    
    # Frontend
    lines.append("    subgraph Frontend[Frontend]")
    lines.append("        front8[front8<br/>Next.js]")
    lines.append("    end")
    
    # Service dependencies
    for service, deps in SERVICE_DEPS.items():
        service_id = service.replace('-', '_')
        for dep in deps:
            dep_id = dep.replace('-', '_')
            lines.append(f"    {service_id} -->|HTTP| {dep_id}")
    
    # External dependencies
    for service, externals in EXTERNAL_DEPS.items():
        service_id = service.replace('-', '_')
        for ext in externals:
            ext_id = ext.replace(' ', '_').replace('-', '_')
            lines.append(f"    {service_id} -->|uses| {ext_id}")
    
    # Shared dependencies
    for shared, services in SHARED_DEPS.items():
        shared_id = shared.replace('-', '_')
        for service in services:
            service_id = service.replace('-', '_')
            lines.append(f"    {service_id} -->|imports| {shared_id}")
    
    return '\n'.join(lines)

def generate_dot():
    """Generate Graphviz DOT format."""
    lines = ["digraph meme9 {"]
    lines.append("    rankdir=LR;")
    lines.append("    node [shape=box];")
    lines.append("")
    
    # Service dependencies
    for service, deps in SERVICE_DEPS.items():
        service_id = service.replace('-', '_')
        for dep in deps:
            dep_id = dep.replace('-', '_')
            lines.append(f'    "{service}" -> "{dep}" [label="HTTP"];')
    
    # External dependencies
    for service, externals in EXTERNAL_DEPS.items():
        for ext in externals:
            lines.append(f'    "{service}" -> "{ext}" [label="uses", style=dashed];')
    
    # Shared dependencies
    for shared, services in SHARED_DEPS.items():
        for service in services:
            lines.append(f'    "{service}" -> "{shared}" [label="imports", style=dotted];')
    
    lines.append("}")
    return '\n'.join(lines)

def generate_text():
    """Generate text representation."""
    lines = ["Dependency Graph for meme9"]
    lines.append("=" * 50)
    lines.append("")
    
    lines.append("Service Dependencies:")
    lines.append("-" * 30)
    for service, deps in sorted(SERVICE_DEPS.items()):
        if deps:
            deps_str = ", ".join(deps)
            lines.append(f"  {service:25} -> {deps_str}")
        else:
            lines.append(f"  {service:25} -> (no dependencies)")
    
    lines.append("")
    lines.append("External Dependencies:")
    lines.append("-" * 30)
    for service, externals in sorted(EXTERNAL_DEPS.items()):
        externals_str = ", ".join(externals)
        lines.append(f"  {service:25} -> {externals_str}")
    
    lines.append("")
    lines.append("Shared Modules:")
    lines.append("-" * 30)
    for shared, services in SHARED_DEPS.items():
        services_str = ", ".join(services)
        lines.append(f"  {shared:25} <- {services_str}")
    
    return '\n'.join(lines)

def main():
    output_dir = os.path.dirname(os.path.abspath(__file__))
    project_root = os.path.dirname(output_dir)
    
    # Generate Mermaid
    mermaid_content = generate_mermaid()
    mermaid_path = os.path.join(project_root, "DEPENDENCY_GRAPH.md")
    with open(mermaid_path, 'w') as f:
        f.write("# Dependency Graph\n\n")
        f.write("```mermaid\n")
        f.write(mermaid_content)
        f.write("\n```\n")
    print(f"✓ Generated Mermaid diagram: {mermaid_path}")
    
    # Generate DOT
    dot_content = generate_dot()
    dot_path = os.path.join(project_root, "dependency_graph.dot")
    with open(dot_path, 'w') as f:
        f.write(dot_content)
    print(f"✓ Generated DOT file: {dot_path}")
    print(f"  Render with: dot -Tpng dependency_graph.dot -o dependency_graph.png")
    
    # Generate text
    text_content = generate_text()
    text_path = os.path.join(project_root, "DEPENDENCY_GRAPH.txt")
    with open(text_path, 'w') as f:
        f.write(text_content)
    print(f"✓ Generated text representation: {text_path}")
    
    print("\n" + "=" * 50)
    print("Dependency Summary:")
    print("=" * 50)
    print(text_content)

if __name__ == "__main__":
    main()

