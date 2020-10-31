import React from "react";
import * as schema from "../../schema/login";

interface PostComposerProps {
    data: schema.ComposerRenderer;
}

interface PostComposerState {
    text: string;
}

export class PostComposer extends React.Component<PostComposerProps, PostComposerState> {
    state: PostComposerState = {
        text: '',
    };

    onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        this.setState({text: e.target.value});
    };

    onSubmit = () => {
        //@ts-ignore
        const params: schema.AnyRequest = {
            addPostRequest: {
                text: this.state.text,
            }
        };

        fetch("/api", {
            method: 'POST',
            body: JSON.stringify(params),
        }).then(r => r.json()).then(r => {

        });
    };

    render() {
        const data = this.props.data;

        return (
            <div>
                Напишите свой пост здесь:<br/>
                {data.welcomeText}<br/>
                <textarea onChange={this.onChange}></textarea><br/>
                <button onClick={this.onSubmit}>Отправить</button>
            </div>
        );
    }
}