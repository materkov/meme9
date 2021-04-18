export function resolveRoute(route: string): Promise<[string, any]> {
    return new Promise(((resolve, reject) => {
        fetch("http://127.0.0.1:8000/router?url=" + encodeURIComponent(route))
            .then(r => r.json())
            .then(r => {
                resolve([r.component, r.data]);
            })
    }));
}
