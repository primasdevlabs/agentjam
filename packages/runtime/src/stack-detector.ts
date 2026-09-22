/**
 * @agentjam/runtime — Stack Detector
 *
 * Inspects project manifests (package.json, composer.json, pyproject.toml, go.mod, etc.)
 * to automatically detect stack technology profiles.
 */

import fs from 'node:fs';
import path from 'node:path';

export interface DetectedStack {
  id: string;
  name: string;
  manifestFile: string;
  frameworks: string[];
}

export class StackDetector {
  static detectStack(projectRoot: string): DetectedStack | undefined {
    // 1. Node / Web Stack
    const pkgPath = path.join(projectRoot, 'package.json');
    if (fs.existsSync(pkgPath)) {
      try {
        const pkg = JSON.parse(fs.readFileSync(pkgPath, 'utf-8'));
        const deps = { ...pkg.dependencies, ...pkg.devDependencies };
        const frameworks: string[] = [];

        if (deps.next) frameworks.push('nextjs');
        if (deps.react) frameworks.push('react');
        if (deps.vue) frameworks.push('vue');
        if (deps.svelte) frameworks.push('svelte');
        if (deps.express) frameworks.push('express');
        if (deps.nestjs || deps['@nestjs/core']) frameworks.push('nestjs');

        if (frameworks.includes('nextjs')) {
          return { id: 'nextjs-fullstack', name: 'Next.js Fullstack', manifestFile: 'package.json', frameworks };
        }
        if (frameworks.includes('react')) {
          return { id: 'react-spa', name: 'React SPA', manifestFile: 'package.json', frameworks };
        }
        return { id: 'node-typescript', name: 'Node.js / TypeScript', manifestFile: 'package.json', frameworks };
      } catch {
        // Fallthrough
      }
    }

    // 2. PHP Stack
    const composerPath = path.join(projectRoot, 'composer.json');
    if (fs.existsSync(composerPath)) {
      try {
        const composer = JSON.parse(fs.readFileSync(composerPath, 'utf-8'));
        const reqs = { ...composer.require, ...composer['require-dev'] };
        const frameworks: string[] = [];

        if (reqs['laravel/framework']) frameworks.push('laravel');
        if (reqs['symfony/symfony'] || reqs['symfony/flex']) frameworks.push('symfony');

        if (frameworks.includes('laravel')) {
          return { id: 'laravel-monolith', name: 'Laravel Modular Monolith', manifestFile: 'composer.json', frameworks };
        }
        return { id: 'php-standard', name: 'PHP Standard', manifestFile: 'composer.json', frameworks };
      } catch {
        // Fallthrough
      }
    }

    // 3. Python Stack
    const pyprojectPath = path.join(projectRoot, 'pyproject.toml');
    const reqTxtPath = path.join(projectRoot, 'requirements.txt');
    if (fs.existsSync(pyprojectPath) || fs.existsSync(reqTxtPath)) {
      return { id: 'python-fastapi', name: 'Python / FastAPI', manifestFile: fs.existsSync(pyprojectPath) ? 'pyproject.toml' : 'requirements.txt', frameworks: ['fastapi', 'python'] };
    }

    // 4. Go Stack
    const goModPath = path.join(projectRoot, 'go.mod');
    if (fs.existsSync(goModPath)) {
      return { id: 'go-microservice', name: 'Go Microservices', manifestFile: 'go.mod', frameworks: ['go'] };
    }

    return undefined;
  }
}
