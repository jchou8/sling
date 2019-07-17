import * as React from 'react'
import { createShallow, createMount } from '@material-ui/core/test-utils'
import SideBarHeader from './SideBarHeader'
import { User } from '../../types'

const logOutFn = jest.fn()

describe('<SideBarHeader />', () => {
  let shallow: Function
  let mount: Function
  beforeAll(() => {
    shallow = createShallow()
    mount = createMount()
  })

  it('should render properly', () => {
    const component = shallow(<SideBarHeader curUser={null} logOut={logOutFn} />)
    expect(component).toMatchSnapshot()
  })

  it('should display the user\'s name', () => {
    const user: User = {
      username: 'Jeff',
      id: 1,
      jwtToken: null
    }

    const component = shallow(<SideBarHeader curUser={user} logOut={logOutFn} />)
    expect(component.containsMatchingElement(<label>Jeff's sling</label>)).toBe(true)
  })

  it('should log the user out upon pressing the log out button', () => {
    const component = mount(<SideBarHeader curUser={null} logOut={logOutFn} />)
    component.find('button').simulate('click')
    expect(logOutFn).toHaveBeenCalled()
  })
})