import {Action, CreatePost, CreatePostReq, CreatePostResp, Navigate, RootRenderer, WriteLog} from "./components/types";

export class State {
    rootRenderer: RootRenderer = {
        label: "Test label",
        command: {
            navigate: {
                url: "/test"
            }
        },
        header: {
            userName: "",
        },
        composer: {
            placeholder: "Напишите что-нибудь..."
        }
    };
    url: string = "";
    callback?: (r: RootRenderer) => void = undefined;

    public do(action: Action) {
        console.log('before', this);

        if (action.writeLog) {
            this.writeLog(action.writeLog);
        } else if (action.navigate) {
            this.navigate(action.navigate);
        } else if (action.createPost) {
            this.createPost(action.createPost);
        } else {
            console.error("Unknown action: ", action);
            return;
        }

        console.log('after', this);
        this.callback && this.callback(this.rootRenderer);
    }

    public subscribe(callback: (r: RootRenderer) => void) {
        this.callback = callback;
    }

    private navigate(params: Navigate) {
        this.url = params.url;
        this.rootRenderer.label = params.url;
        window.history.pushState(null, 'meme9', params.url);
    }

    private writeLog(params: WriteLog) {
        console.log('Write Log!', params.message);
    }

    private createPost(params: CreatePost): Promise<CreatePostResp> {
        return new Promise<CreatePostResp>((resolve, reject) => {
            const req: CreatePostReq = {
                text: params.text,
            };

            fetch("http://localhost:8001/CreatePost", {
                method: "POST",
                body: JSON.stringify(req),
            })
                .then(r => r.json())
                .then(r => {
                    resolve(r)
                })
        })
    }
}


//export const state = new State();
///** @ts-ignore */
//window.state = state;
