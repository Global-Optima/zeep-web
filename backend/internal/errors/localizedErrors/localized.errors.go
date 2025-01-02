package localizedErrors

type LocalizedErrorInterface interface {
	GetEn() string
	GetRu() string
	GetKk() string
}

type LocalizedError struct {
	en string
	ru string
	kk string
}

func (e *LocalizedError) GetEn() string {
	return e.en
}

func (e *LocalizedError) GetRu() string {
	return e.ru
}

func (e *LocalizedError) GetKk() string {
	return e.kk
}

func NewLocalizedError(en, ru, kk string) *LocalizedError {
	return &LocalizedError{en: en, ru: ru, kk: kk}
}
