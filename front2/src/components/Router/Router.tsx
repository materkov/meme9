import React from "react";
import {GlobalContext, GlobalStore, GlobalStoreContext} from "../../Context";
import {Header} from "../Header/Header";
import * as schema from "../../api/renderer";
import {UniversalRenderer} from "../UniversalRenderer/UniversalRenderer";
import {API} from "../../Api";
import {ToggleLikeRequest_Action} from "../../api/posts";

interface State {
    data?: schema.UniversalRenderer;
    error?: boolean;
}

export class Router extends React.Component<{}, State> {
    state: State = {};
    globalStore: GlobalStore;

    constructor(props: any) {
        super(props);
        this.globalStore = {
            // TODO mutations here. Think about it
            togglePostLike: (postId: string) => {
                if (this.state.data?.feedRenderer) {
                    for (let post of this.state.data.feedRenderer.posts) {
                        if (post.id == postId && post.canLike) {
                            let action: ToggleLikeRequest_Action;

                            if (post.isLiked) {
                                post.isLiked = false;
                                post.likesCount--;
                                action = ToggleLikeRequest_Action.UNLIKE;
                            } else {
                                post.isLiked = true;
                                post.likesCount++;
                                action = ToggleLikeRequest_Action.LIKE;
                            }

                            API.Posts_ToggleLike({
                                action: action,
                                postId: post.id,
                            }).then(r => {
                                post.likesCount = r.likesCount;
                                this.setState({data: this.state.data}); // TODO think about this hack
                            }).catch(e => {
                                console.error(e);
                            });
                            break;
                        }
                    }
                }
            }
        }
    }


    componentDidMount() {
        this.navigate(window.location.pathname);
    }

    navigate = (route: string) => {
        window.history.pushState(null, "meme", route);

        API.Utils_ResolveRoute({url: route})
            .then(data => this.setState({data: schema.UniversalRenderer.fromJSON(data), error: undefined}))
            .catch(() => this.setState({data: undefined, error: true}))
    }

    render() {
        return (
            <GlobalContext.Provider value={this.navigate}>
                <GlobalStoreContext.Provider value={this.globalStore}>
                    <Header/>
                    {this.state.data && <UniversalRenderer data={this.state.data}/>}
                    {this.state.error && <div>Ошибка!</div>}
                </GlobalStoreContext.Provider>
            </GlobalContext.Provider>
        )
    }
}