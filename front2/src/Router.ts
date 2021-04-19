export function resolveRoute(route: string): Promise<[string, any]> {
    return new Promise(((resolve, reject) => {
        fetch("http://localhost:8000/router?url=" + encodeURIComponent(route), {
            credentials: 'include',
        })
            .then(r => r.json())
            .then(r => {
                resolve([r.component, r.data]);
            })
    }));
}
