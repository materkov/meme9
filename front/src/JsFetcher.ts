let globalLoaded: { [file: string]: boolean } = {};

export function fetchJs(js: string[]): Promise<null> {
    js = js || [];

    return new Promise<null>((resolve, reject) => {
        let neededJs = [];
        for (let file of js) {
            if (!globalLoaded[file]) {
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
                    globalLoaded[file] = true;
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
