import React from "react";
import styles from "./Composer.module.css";
import classNames from "classnames";
import {Link} from "../Link/Link";
import {GlobalStoreContext} from "../../Context";
import {API} from "../../Api";
import {Store} from "../../Store";

interface State {
    text: string;
    isTextFocused?: boolean;
    addedPostUrl?: string;
    uploadedPhoto?: string;
    uploadedPhotoID: string;
    isUploading: boolean;
}

export class Composer extends React.Component<{}, State> {
    static contextType = GlobalStoreContext;

    state: State = {
        text: '',
        isUploading: false,
        uploadedPhotoID: "",
    };

    onChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
        this.setState({text: e.currentTarget.value});
    };

    onSubmit = () => {
        const store = this.context as Store;
        store.addPost({
            text: this.state.text,
            photoId: this.state.uploadedPhotoID,
        })
            .then(r => this.setState({addedPostUrl: r.postUrl}))
            .catch(() => console.error("Error saving post"))

        this.setState({text: ''});
    };

    onTextFocus = () => {
        this.setState({isTextFocused: true});
    }

    onTextBlur = () => {
        this.setState({isTextFocused: false});
    }

    onUploadPhoto = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (!e.target.files || !e.target.files[0]) {
            return;
        }
        if (this.state.isUploading) {
            return;
        }

        this.setState({isUploading: true});

        API.Upload(e.target.files[0])
            .then((r: any) => {
                this.setState({
                    uploadedPhoto: r.url,
                    uploadedPhotoID: r.id,
                    isUploading: false,
                })
            })
            .catch(console.error)
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
                })}>

                    <div className={classNames({
                        [styles.UploadPhoto__hidden]: this.state.uploadedPhoto,
                        [styles.UploadPhotoPreview]: !this.state.uploadedPhoto,
                    })}>
                        <input type="file" onChange={this.onUploadPhoto}/>
                    </div>

                    <button className={styles.SubmitBtn} onClick={this.onSubmit}>Опубликовать</button>
                </div>

                {this.state.uploadedPhoto &&
                <img className={styles.UploadPhotoPreview} src={this.state.uploadedPhoto}/>
                }

                {this.state.addedPostUrl &&
                <Link href={this.state.addedPostUrl}>Пост добавлен</Link>
                }
            </div>
        );
    }
}
