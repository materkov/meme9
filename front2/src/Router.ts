import {UniversalRenderer} from "./api/renderer";

export function resolveRoute(route: string): Promise<UniversalRenderer> {
    return new Promise(((resolve, reject) => {
        fetch("http://localhost:8000" + route, {
            credentials: 'include',
        })
            .then(r => r.json())
            .then(r => {
                resolve(UniversalRenderer.fromJSON(r.data));
            })
    }));
}
