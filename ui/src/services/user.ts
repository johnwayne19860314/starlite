import request from '@/utils/request';

//登录
export const userLogin = (params: Record<string, unknown>): Promise<API.LoginResponse> => (
  request.post('/api/starlite/v1/user/login', params)
);

export const userList = (params: Record<string, unknown>): Promise<API.UsersListResponse> => (
  request.post('/api/starlite/v1/user/list', params)
);



//获取用户信息
export const retrieveUserInfo = (): Promise<API.UserInfoResponse> => (
  request.get('/api/starlite/v1/user')
);

export const addUser = (params: Record<string, API.UserInfo>): Promise<API.UserInfoResponse> => (
  request.post('/api/starlite/v1/user',params)
);
export const updateUser = (params: Record<string, API.UserInfo>): Promise<API.UserInfoResponse> => (
  request.put('/api/starlite/v1/user',params)
);

export const delUser = (id:  number): Promise<API.DelUserResponse> => (
  request.delete(`/api/starlite/v1/user/${id}`)
);
//获取用户权限
export const retrieveUserAuthority = (): Promise<API.UserAuthorityResponse> => (
  request.get('/api/starlite/v1/user/authority')
);

//获取菜单数据
export const retrieveMenuData = (): Promise<API.MenuDataResponse> => (
  request.get('/api/starlite/v1/user/menu')
);

//获取用户信息和权限以及菜单
export const retrieveUserInfoAuthorityMenu = (): Promise<API.UserInfoAuthMenuResponse> => (
  Promise.all([
    retrieveUserInfo(),
    retrieveUserAuthority(),
    retrieveMenuData(),
  ])
);

//获取用户权限以及菜单
export const retrieveUserAuthorityMenu = (): Promise<API.UserAuthMenuResponse> => (
  Promise.all([
    retrieveUserAuthority(),
    retrieveMenuData(),
  ])
);

//登出
export const userLogout = (): Promise<API.LogoutResponse> => (
  request.post('/api/starlite/v1/user/logout')
);

//获取验证码
export const retrieveCaptcha = (params: Record<string, string>): Promise<API.CaptchaResponse> => (
  request.get('/api/starlite/v1/user/captcha', { params })
);
