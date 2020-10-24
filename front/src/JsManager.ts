let loaded: { [file: string]: boolean } = {};

if (window.InitJsBundles) {
    for (let js of window.InitJsBundles) {
        loaded[js] = true
    }
}

export function loadJs(js: string[]): Promise<null> {
    return new Promise<null>((resolve, reject) => {
        let neededJs = [];
        for (let file of js) {
            if (!loaded[file]) {
                neededJs.push(file);
            }
        }

        if (neededJs.length === 0) {
            resolve();
        }

        let promises = [];
        for (let file of neededJs) {
            promises.push(new Promise<null>(((resolve, reject) => {
                const script = document.createElement('script');
                script.src = file;
                script.onload = () => {
                    loaded[file] = true;
                    resolve();
                };

                document.body.appendChild(script);
            })))
        }

        Promise.all(promises).then(() => {
            resolve();
        })
    });
}
