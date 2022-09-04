import React from "react";

const registry: { [key: string]: React.Component } = {};

export function RegisterComponent(name: string, component: React.Component) {
    registry[name] = component;
}

export function GetComponent(name: string): React.Component | undefined {
    return registry[name];
}
