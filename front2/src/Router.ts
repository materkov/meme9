import {UniversalRenderer} from "./api/renderer";

export function resolveRoute(route: string): Promise<UniversalRenderer> {
    return new Promise((resolve, reject) => {
        fetch("http://localhost:8000" + route, {
            credentials: 'include',
        })
            .then(r => {
                if (r.status !== 200) {
                    reject();
                    return;
                }

                return r.json()
            })
            .then(r => {
                resolve(UniversalRenderer.fromJSON(r.data))
            })
            .catch(() => reject())
    });
}
