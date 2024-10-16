package makosh

//var ErrNoUpdaterFunction = errors.New("no updater function")
//
//type DefaultResolver struct {
//	getAddressesRequest *http.Request
//
//	m sync.RWMutex
//
//	doUpdateAddresses makosh_resolver.UpdateAddresses
//
//	log logrus.StdLogger
//
//	overrides []string
//}
//
//func (r *DefaultResolver) UpdateAddresses(addresses ...string) error {
//	r.m.Lock()
//	defer r.m.Unlock()
//
//	if r.doUpdateAddresses == nil {
//		return errors.Wrap(ErrNoUpdaterFunction)
//	}
//
//	err := r.doUpdateAddresses(addresses)
//	if err != nil {
//		return errors.Wrap(err, "error updating addresses")
//	}
//
//	return nil
//}
//
//func (r *DefaultResolver) AddUpdater(updater makosh_resolver.UpdateAddresses) {
//	r.m.Lock()
//	r.doUpdateAddresses = updater
//	r.m.Unlock()
//}
//
//func (r *DefaultResolver) Resolve() error {
//	if len(r.overrides) != 0 {
//		return r.doUpdateAddresses(r.overrides)
//	}
//
//	httpResp, err := http.DefaultClient.Do(r.getAddressesRequest)
//	if err != nil {
//		return errors.Wrap(err, "error getting ")
//	}
//
//	body, err := io.ReadAll(httpResp.Body)
//	if err != nil {
//		return errors.Wrap(err, "error reading body")
//	}
//
//	if httpResp.StatusCode != http.StatusOK {
//		return errors.New(string(body))
//	}
//
//	var endpointsResponse makosh_be.ListEndpoints_Response
//	err = json.Unmarshal(body, &endpointsResponse)
//	if err != nil {
//		return errors.Wrap(err, "error parsing list endpoints response")
//	}
//
//	err = r.doUpdateAddresses(endpointsResponse.Urls)
//	if err != nil {
//		return errors.Wrap(err, "error updating connection state")
//	}
//
//	return nil
//}
//
//func (r *DefaultResolver) Close() {}
