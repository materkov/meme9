import React from "react";
import * as schema from "../../schema/api";
import {fetchData} from "../../DataFetcher";
import {Link} from "../Link/Link";

interface PostComposerProps {
    data: schema.ComposerRenderer;
}

interface PostComposerState {
    text: string;
    success?: schema.AddPostRenderer;
}

export class PostComposer extends React.Component<PostComposerProps, PostComposerState> {
    state: PostComposerState = {
        text: '',
    };

    onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        this.setState({text: e.target.value});
    };

    onSubmit = () => {
        const params: schema.AddPostRequest = {
            text: this.state.text,
        };

        fetchData<schema.AddPostRenderer>('meme.API.AddPost', JSON.stringify(params)).then(r => {
            this.setState({success: r});
        }).catch(() => {

        })
    };

    render() {
        const data = this.props.data;

        return (
            <div>
                {data.welcomeText}<br/>
                <textarea onChange={this.onChange}/><br/>
                <button onClick={this.onSubmit}>{data.sendText}</button>

                {this.state.success &&
                <div>
                    {this.state.success.successText}
                    <Link href={this.state.success.postUrl}>Перейти</Link>
                </div>
                }
            </div>
        );
    }
}