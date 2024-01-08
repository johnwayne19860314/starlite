
export const avatarVal = 'https://gw.alipayobjects.com/zos/antfincdn/XAosXuNZyF/BiazfanxmamNRoxxVxka.png'

export const initUserInfo : API.UserInfo = {
  id:0,
  name: 'Serati Ma',
  avatar: avatarVal,
  //userid: '00000001',
  email: 'antdesign@alipay.com',
  //signature: '海纳百川，有容乃大',
  role: 2,
}

export const userInfo : API.UserInfoResponse = {
  data: initUserInfo,
  code : 0,
  message : "success"
  }
    
  export const internalAuthList = [
    '/',
    '/about',
  ]
  export const powerAuthList = [
    ...internalAuthList,
    '/entry/detail',
    '/entry/category',

  ]
  export const adminAuthList = [
    ...powerAuthList,
    '/users',

  ]
  

  export const makeUserAuthority  = (userAuthority:string[]) => {
   
    let res :API.UserAuthorityResponse = {
      data:{
        authority: [
          ...userAuthority
        ]
      },
      code:0,
      message:"success"
    }
    return res
}

  export const internalMenu :API.MenuData= [{
    id: 1,
    key: '1',
    path: '/',
    redirect: '',
    label: '首页'
  }]

  export const powerMenu :API.MenuData= [
    ...internalMenu,
    {
      id: 2,
      key: '2',
      label: '关于',
      path: '/about',
      redirect: '/about',
      
    },
    {
      id: 3,
      key: '3',
      label: '原料',
      path: '/entry',
      redirect: '/entry/detail',
      children: [
        
        {
          id: 31,
          key: '3-1',
          path: '/entry/category',
          redirect: '',
          label: '种类',
          pid: 3
        },
        {
          id: 32,
          key: '3-2',
          path: '/entry/detail',
          redirect: '',
          label: '原料',
          pid: 3
        },
        
      ]
    }
  ]
  export const adminMenu :API.MenuData= [
    ...powerMenu,
    {
      id: 4,
      key: '4',
      label: '用户',
      path: '/users',
      redirect: '/users',
      
    },
  ]
export const makeUserMenu  = (userMenu:API.MenuData) => {
   
    let res :API.MenuDataResponse = {
      data:userMenu,
      code:0,
      message:"success"
    }
    return res
}

  // export const externalUserMenu :API.MenuDataResponse= {
  //   data : 
  //     externalMenu,
  //     code : 0,
  //     message : "success"
  // }
  // export const powerUserMenu :API.MenuDataResponse= {
  //   data : 
  //     internalMenu,
  //     code : 0,
  //     message : "success"
  // }
  // export const adminMenu :API.MenuDataResponse= {
  //   data : 
  //     [
  //       {
  //         id: 1,
  //         key: '1',
  //         path: '/',
  //         redirect: '',
  //         label: '首页'
  //       },
  //       {
  //         id: 2,
  //         key: '2',
  //         label: '关于',
  //         path: '/about',
  //         redirect: '/about',
         
  //       },
  //       {
  //         id: 3,
  //         key: '3',
  //         label: '原料',
  //         path: '/entry',
  //         redirect: '/entry/detail',
  //         children: [
            
  //           {
  //             id: 32,
  //             key: '3-2',
  //             path: '/entry/detail',
  //             redirect: '',
  //             label: '原料',
  //             pid: 3
  //           },
            
  //         ]
  //       }
  //     ],
  //     code : 0,
  //     message : "success"
    
  // }

  export const userLogoutData : API.LogoutResponse = {
    data : {},
    code : 0,
    message : "success"
  }