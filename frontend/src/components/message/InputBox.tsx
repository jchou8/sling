import * as React from 'react';
import './styles.scss';
import {
    TextField,
    IconButton,
    InputAdornment
} from '@material-ui/core'
import Send from '@material-ui/icons/Send'

type InputBoxProps = {
    sendMessage: Function
    enabled: boolean
}

type InputBoxState = {
    input: string
}

class InputBox extends React.Component<InputBoxProps, InputBoxState> {
    constructor(props: InputBoxProps) {
        super(props);
        this.state = {
            input: "",
        }
    }
    handleInputChange = (value: string) => {
        this.setState({ input: value });
    }

    handleSubmit = (e: React.FormEvent) => {
        e.preventDefault();
        this.props.sendMessage(this.state.input)
        this.setState({
            input: "",
        });
    }

    render() {
        return (
            <div className="InputBox">
                <form className="input-form" onSubmit={this.handleSubmit}>
                    <TextField
                        className="input-box"
                        id="input-box"
                        value={this.state.input}
                        onChange={e => this.handleInputChange(e.currentTarget.value)}
                        margin="dense"
                        fullWidth
                        variant="filled"
                        disabled={!this.props.enabled}
                        InputProps={{
                            endAdornment: (
                                <InputAdornment position="end">
                                    <IconButton
                                        id="input-submit"
                                        className="input-submit"
                                        type="submit"
                                        value="Submit"
                                        disabled={!this.props.enabled}
                                    >
                                        <Send />
                                    </IconButton>
                                </InputAdornment>
                            )
                        }}
                    />
                </form>
            </div>
        );
    }
}

export default InputBox;