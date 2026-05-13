package logctx

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
	"testing"
	"time"
)

type mockSlogHandler struct {
	mock.Mock
}

func (m *mockSlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return m.Called(ctx, level).Bool(0)
}

func (m *mockSlogHandler) Handle(ctx context.Context, r slog.Record) error {
	return m.Called(ctx, r).Error(0)
}

func (m *mockSlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return m.Called(attrs).Get(0).(slog.Handler)
}

func (m *mockSlogHandler) WithGroup(name string) slog.Handler {
	return m.Called(name).Get(0).(slog.Handler)
}

func matchRecordAttrs(attrs []slog.Attr) func(slog.Record) bool {
	return func(r slog.Record) bool {
		var got []slog.Attr
		r.Attrs(func(a slog.Attr) bool {
			got = append(got, a)
			return true
		})
		return assert.ObjectsAreEqual(got, attrs)
	}
}

func TestNewSlogHandler(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "Returns non-nil slog.Handler"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inner := new(mockSlogHandler)
			got := NewSlogHandler(inner)
			assert.NotNilf(t, got, "NewSlogHandler(%v)", inner)
		})
	}
}

func TestSlogHandler_Enabled(t *testing.T) {
	type args struct {
		ctx   context.Context
		level slog.Level
	}
	type mockOpts struct {
		handler func(m *mockSlogHandler)
	}
	tests := []struct {
		name     string
		args     args
		mockOpts mockOpts
		want     bool
	}{
		{
			name: "Inner handler returns true, Enabled returns true",
			args: args{ctx: context.Background(), level: slog.LevelInfo},
			mockOpts: mockOpts{handler: func(m *mockSlogHandler) {
				m.On("Enabled", mock.Anything, mock.Anything).Return(true)
			}},
			want: true,
		},
		{
			name: "Inner handler returns false, Enabled returns false",
			args: args{ctx: context.Background(), level: slog.LevelDebug},
			mockOpts: mockOpts{handler: func(m *mockSlogHandler) {
				m.On("Enabled", mock.Anything, mock.Anything).Return(false)
			}},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inner := new(mockSlogHandler)
			if tt.mockOpts.handler != nil {
				tt.mockOpts.handler(inner)
			}
			h := NewSlogHandler(inner)
			got := h.Enabled(tt.args.ctx, tt.args.level)
			assert.Equalf(t, tt.want, got, "Enabled(%v, %v)", tt.args.ctx, tt.args.level)
			inner.AssertExpectations(t)
		})
	}
}

func TestSlogHandler_Handle(t *testing.T) {
	type args struct {
		ctx    context.Context
		record slog.Record
	}
	type mockOpts struct {
		handler func(m *mockSlogHandler)
	}
	tests := []struct {
		name     string
		args     args
		mockOpts mockOpts
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "Context has no LogCtx, inner Handle called with no extra attrs",
			args: args{
				ctx:    context.Background(),
				record: slog.NewRecord(time.Time{}, slog.LevelInfo, "test message", 0),
			},
			mockOpts: mockOpts{handler: func(m *mockSlogHandler) {
				m.On("Handle", mock.Anything, mock.MatchedBy(matchRecordAttrs(nil))).Return(nil)
			}},
			wantErr: assert.NoError,
		},
		{
			name: "Context has LogCtx fields, fields added to record before inner Handle",
			args: args{
				ctx:    Add(context.Background(), String("key", "value"), Int("count", 42)),
				record: slog.NewRecord(time.Time{}, slog.LevelInfo, "test message", 0),
			},
			mockOpts: mockOpts{handler: func(m *mockSlogHandler) {
				m.On("Handle", mock.Anything, mock.MatchedBy(matchRecordAttrs(
					[]slog.Attr{slog.String("key", "value"), slog.Int("count", 42)},
				))).Return(nil)
			}},
			wantErr: assert.NoError,
		},
		{
			name: "Inner Handle returns error, error is returned",
			args: args{
				ctx:    context.Background(),
				record: slog.NewRecord(time.Time{}, slog.LevelInfo, "test message", 0),
			},
			mockOpts: mockOpts{handler: func(m *mockSlogHandler) {
				m.On("Handle", mock.Anything, mock.MatchedBy(matchRecordAttrs(nil))).Return(errors.New("handle error"))
			}},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inner := new(mockSlogHandler)
			if tt.mockOpts.handler != nil {
				tt.mockOpts.handler(inner)
			}
			h := NewSlogHandler(inner)
			err := h.Handle(tt.args.ctx, tt.args.record)
			tt.wantErr(t, err, "Handle(%v, %v)", tt.args.ctx, tt.args.record)
			inner.AssertExpectations(t)
		})
	}
}

func TestSlogHandler_WithAttrs(t *testing.T) {
	type args struct {
		attrs []slog.Attr
	}
	type mockOpts struct {
		handler func(m *mockSlogHandler)
	}
	type testCase struct {
		name      string
		args      args
		mockOpts  mockOpts
		wantInner *mockSlogHandler
	}

	innerResult1 := new(mockSlogHandler)
	innerResult2 := new(mockSlogHandler)

	tests := []testCase{
		{
			name:      "Returns slogHandler wrapping inner.WithAttrs result",
			args:      args{attrs: []slog.Attr{slog.String("key", "value")}},
			wantInner: innerResult1,
			mockOpts: mockOpts{handler: func(m *mockSlogHandler) {
				m.On("WithAttrs", []slog.Attr{slog.String("key", "value")}).Return(slog.Handler(innerResult1))
			}},
		},
		{
			name:      "Returns slogHandler wrapping inner.WithAttrs with empty attrs",
			args:      args{attrs: []slog.Attr{}},
			wantInner: innerResult2,
			mockOpts: mockOpts{handler: func(m *mockSlogHandler) {
				m.On("WithAttrs", []slog.Attr{}).Return(slog.Handler(innerResult2))
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inner := new(mockSlogHandler)
			if tt.mockOpts.handler != nil {
				tt.mockOpts.handler(inner)
			}
			h := NewSlogHandler(inner)
			got := h.WithAttrs(tt.args.attrs)
			assert.NotNilf(t, got, "WithAttrs(%v)", tt.args.attrs)
			gotSlogHandler, ok := got.(*slogHandler)
			assert.Truef(t, ok, "WithAttrs(%v) returns *slogHandler", tt.args.attrs)
			if ok {
				assert.Equalf(t, slog.Handler(tt.wantInner), gotSlogHandler.inner, "WithAttrs(%v) inner", tt.args.attrs)
			}
			inner.AssertExpectations(t)
		})
	}
}

func TestSlogHandler_WithGroup(t *testing.T) {
	type args struct {
		name string
	}
	type mockOpts struct {
		handler func(m *mockSlogHandler)
	}
	type testCase struct {
		name      string
		args      args
		mockOpts  mockOpts
		wantInner *mockSlogHandler
	}

	innerResult1 := new(mockSlogHandler)
	innerResult2 := new(mockSlogHandler)

	tests := []testCase{
		{
			name:      "Returns slogHandler wrapping inner.WithGroup result",
			args:      args{name: "test-group"},
			wantInner: innerResult1,
			mockOpts: mockOpts{handler: func(m *mockSlogHandler) {
				m.On("WithGroup", "test-group").Return(slog.Handler(innerResult1))
			}},
		},
		{
			name:      "Returns slogHandler wrapping inner.WithGroup with empty name",
			args:      args{name: ""},
			wantInner: innerResult2,
			mockOpts: mockOpts{handler: func(m *mockSlogHandler) {
				m.On("WithGroup", "").Return(slog.Handler(innerResult2))
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inner := new(mockSlogHandler)
			if tt.mockOpts.handler != nil {
				tt.mockOpts.handler(inner)
			}
			h := NewSlogHandler(inner)
			got := h.WithGroup(tt.args.name)
			assert.NotNilf(t, got, "WithGroup(%v)", tt.args.name)
			gotSlogHandler, ok := got.(*slogHandler)
			assert.Truef(t, ok, "WithGroup(%v) returns *slogHandler", tt.args.name)
			if ok {
				assert.Equalf(t, slog.Handler(tt.wantInner), gotSlogHandler.inner, "WithGroup(%v) inner", tt.args.name)
			}
			inner.AssertExpectations(t)
		})
	}
}
