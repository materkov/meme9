import React from "react";
import styles from "./Composer.module.css";
import classNames from "classnames";
import {Link} from "../Link/Link";
import {GlobalStoreContext} from "../../Context";

interface State {
    text: string;
    isTextFocused?: boolean;
    addedPostUrl?: string;
}

export class Composer extends React.Component<{}, State> {
    static contextType = GlobalStoreContext;

    state: State = {
        text: '',
    };

    onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        this.setState({text: e.currentTarget.value});
    };

    onSubmit = () => {
        this.context.addPost(this.state.text)
            .then((r: string) => this.setState({addedPostUrl: r}))
            .catch(() => console.error("Error saving post"))

        this.setState({text: ''});
    };

    onTextFocus = () => {
        this.setState({isTextFocused: true});
    }

    onTextBlur = () => {
        this.setState({isTextFocused: false});
    }

    render() {
        const expandedState = this.state.text || this.state.isTextFocused;

        return (
            <div>
                <textarea placeholder="Что у Вас нового?" onChange={this.onChange} value={this.state.text}
                          onFocus={this.onTextFocus} onBlur={this.onTextBlur}
                          className={classNames({
                              [styles.Composer]: true,
                              [styles.Composer__hasText]: expandedState,
                          })}
                />

                <div className={classNames({
                    [styles.BottomContainer]: true,
                    [styles.BottomContainer__hidden]: !expandedState,
                })}>
                    <button className={styles.SubmitBtn} onClick={this.onSubmit}>Опубликовать</button>
                </div>

                {this.state.addedPostUrl &&
                <Link href={this.state.addedPostUrl}>Пост добавлен</Link>
                }
            </div>
        );
    }
}
