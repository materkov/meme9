import React from "react";
import * as schema from "../../schema/login";
import {fetchData} from "../../DataFetcher";

interface PostComposerProps {
    data: schema.ComposerRenderer;
}

interface PostComposerState {
    text: string;
    success: boolean;
}

export class PostComposer extends React.Component<PostComposerProps, PostComposerState> {
    state: PostComposerState = {
        text: '',
        success: false,
    };

    onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        this.setState({text: e.target.value});
    };

    onSubmit = () => {
        const params: schema.AddPostRequest = {
            text: this.state.text,
        };

        fetchData<schema.AddPostRenderer>('meme.API.AddPost', JSON.stringify(params)).then(r => {
            this.setState({success: true});
        }).catch(() => {

        })
    };

    render() {
        const data = this.props.data;

        return (
            <div>
                Напишите свой пост здесь:<br/>
                {data.welcomeText}<br/>
                <textarea onChange={this.onChange}></textarea><br/>
                <button onClick={this.onSubmit}>Отправить</button>

                {this.state.success &&
                <div>Пост успешно добавлен</div>
                }
            </div>
        );
    }
}