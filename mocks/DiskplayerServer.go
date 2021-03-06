// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import http "net/http"
import mock "github.com/stretchr/testify/mock"
import oauth2 "golang.org/x/oauth2"
import spotify "github.com/zmb3/spotify"

// DiskplayerServer is an autogenerated mock type for the DiskplayerServer type
type DiskplayerServer struct {
	mock.Mock
}

// Authenticator provides a mock function with given fields:
func (_m *DiskplayerServer) Authenticator() *spotify.Authenticator {
	ret := _m.Called()

	var r0 *spotify.Authenticator
	if rf, ok := ret.Get(0).(func() *spotify.Authenticator); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*spotify.Authenticator)
		}
	}

	return r0
}

// RunCallbackServer provides a mock function with given fields:
func (_m *DiskplayerServer) RunCallbackServer() (*http.Server, error) {
	ret := _m.Called()

	var r0 *http.Server
	if rf, ok := ret.Get(0).(func() *http.Server); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*http.Server)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RunRecordServer provides a mock function with given fields:
func (_m *DiskplayerServer) RunRecordServer() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TokenChannel provides a mock function with given fields:
func (_m *DiskplayerServer) TokenChannel() chan *oauth2.Token {
	ret := _m.Called()

	var r0 chan *oauth2.Token
	if rf, ok := ret.Get(0).(func() chan *oauth2.Token); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan *oauth2.Token)
		}
	}

	return r0
}
