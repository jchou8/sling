import * as React from 'react'
import { createShallow, createMount } from '@material-ui/core/test-utils'
import InputBox from './InputBox'

const sendFn = jest.fn()

describe('<InputBox />', () => {
  let mount: Function
  let shallow: Function
  beforeAll(() => {
    mount = createMount()
    shallow = createShallow()
  })

  it('should render when enabled', () => {
    const component = shallow(<InputBox enabled={true} sendMessage={sendFn} />)
    expect(component).toMatchSnapshot()
  })

  it('should render when disabled', () => {
    const component = shallow(<InputBox enabled={false} sendMessage={sendFn} />)
    expect(component).toMatchSnapshot()
  })

  it('should send a message upon submitting', () => {
    const component = mount(<InputBox enabled={true} sendMessage={sendFn} />)
    component.setState({ input: 'hello' })
    component.find('button#input-submit').simulate('submit')
    expect(sendFn).toHaveBeenCalledWith('hello')
  })
})