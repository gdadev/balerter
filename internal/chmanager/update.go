package manager

//func (m *ChannelsManager) Update(name string, level alert.Level, text string, options *alert.Options) error {
//	m.logger.Debug("alertManager update",
//		zap.String("name", name),
//		zap.Int("level", int(level)),
//		zap.String("text", text),
//		zap.Any("options", options),
//	)
//
//	metrics.SetAlertLevel(name, level)
//
//	a, err := m.storage.Alert().GetOrNew(name)
//	if err != nil {
//		return fmt.Errorf("error get alert %s, %w", name, err)
//	}
//
//	if a.HasLevel(level) {
//		a.Inc()
//
//		if !options.Quiet && options.Repeat > 0 && a.Count()%options.Repeat == 0 {
//			return m.sendMessageFunc(level.String(), name, text, options, m.errs)
//		}
//
//		return nil
//	}
//
//	a.UpdateLevel(level)
//
//	if !options.Quiet {
//		err = m.sendMessageFunc(level.String(), name, text, options, m.errs)
//		if err != nil {
//			return err
//		}
//	}
//
//	err = m.storage.Alert().Store(a)
//	if err != nil {
//		return fmt.Errorf("error release alert, %w", err)
//	}
//
//	return nil
//}
